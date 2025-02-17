package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/golang-jwt/jwt"
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/redis"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
	"github.com/thejixer/jixifood/services/auth/internal/utils"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthLogic struct {
	dbStore    *repository.PostgresStore
	RedisStore *redis.RedisStore
}

// this is the service layer, but since the term service is repeated in different meanings,
// I renamed it to AuthLogic instead of AuthService for clarity
func NewAuthLogic(dbStore *repository.PostgresStore, redisStore *redis.RedisStore) *AuthLogic {
	return &AuthLogic{
		dbStore:    dbStore,
		RedisStore: redisStore,
	}
}

func (l *AuthLogic) GenerateAndStoreOtp(phoneNumber string) (string, error) {

	otp := utils.GenerateOTP()
	err := l.RedisStore.SetOTP(phoneNumber, otp)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (l *AuthLogic) VerifyOTP(phoneNumber, otp string) (bool, error) {
	generatedOTP, err := l.RedisStore.GetOTP(phoneNumber)
	if err != nil {
		return false, apperrors.ErrNotFound
	}

	if generatedOTP != otp {
		return false, apperrors.ErrCodeMismatch
	}

	go l.RedisStore.DelOTP(phoneNumber)

	return true, nil
}

func (l *AuthLogic) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.GetUserByPhoneNumber(ctx, phoneNumber)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			// this behaviour is expected
			return nil, nil
		}
		return nil, apperrors.ErrUnexpected
	}

	return user, nil
}

func (l *AuthLogic) CreateUser(ctx context.Context, phoneNumber, name string, roleID uint64) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.CreateUser(ctx, phoneNumber, name, roleID)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (l *AuthLogic) GetRequester(ctx context.Context) (*models.UserEntity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, apperrors.ErrMissingMetaData
	}
	tokens := md["auth"]
	if len(tokens) == 0 {
		return nil, apperrors.ErrMissingToken
	}
	tokenString := tokens[0]
	token, err := utils.VerifyToken(tokenString)

	if err != nil || !token.Valid {
		return nil, apperrors.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["id"] == nil {
		return nil, apperrors.ErrUnauthorized
	}
	i := claims["id"].(string)
	intId, err := strconv.Atoi(i)
	if err != nil {
		return nil, apperrors.ErrInternal
	}
	userId := uint64(intId)

	userFromCache := l.RedisStore.GetUser(userId)
	if userFromCache != nil {
		return userFromCache, nil
	}

	user, err := l.dbStore.AuthRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	go l.RedisStore.CacheUser(user)

	return user, nil

}

func (l *AuthLogic) ConvertToPBUser(ctx context.Context, user *models.UserEntity) *pb.User {

	role, err := l.dbStore.AuthRepo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		return nil
	}

	return &pb.User{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Status:      GetUserStatus(user.Status),
		Role:        role.Name,
		CreatedAt:   timestamppb.New(user.CreatedAt),
	}
}

func (l *AuthLogic) CheckPermission(ctx context.Context, roleID uint64, permissionName string) bool {

	ok, err := l.dbStore.AuthRepo.CheckPermission(ctx, roleID, permissionName)
	if err != nil || !ok {
		return false
	}

	return true
}

func (l *AuthLogic) ChangeUserRole(ctx context.Context, userID, roleID uint64) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.ChangeUserRole(ctx, userID, roleID)
	if err != nil {
		return nil, fmt.Errorf("error in authLogic.changeUserRole: %w", err)

	}

	return user, nil
}

func (l *AuthLogic) EditProfile(ctx context.Context, userID uint64, name string) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.EditProfile(ctx, userID, name)

	if err != nil {
		return nil, fmt.Errorf("error in authLogic.editProfile: %w", err)
	}

	return user, nil
}

func (l *AuthLogic) QueryUsers(ctx context.Context, text string, page, limit uint64) ([]*pb.User, uint64, bool, error) {
	var (
		userEntities []*models.UserEntity
		count        uint64
		hasNextPage  bool
		roles        []*models.Role
		err1, err2   error
		wg           sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		userEntities, count, hasNextPage, err1 = l.dbStore.AuthRepo.QueryUsers(ctx, text, page, limit)
	}()

	go func() {
		defer wg.Done()
		roles, err2 = l.dbStore.AuthRepo.GetRoles(ctx)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, 0, false, fmt.Errorf("error in authLogic.queryUsers: %w", err1)
	}
	if err2 != nil {
		return nil, 0, false, fmt.Errorf("error in authLogic.getRoles: %w", err2)
	}

	roleCache := make(map[uint64]string)
	for _, role := range roles {
		roleCache[role.ID] = role.Name
	}
	var users []*pb.User
	for _, u := range userEntities {
		user := &pb.User{
			Id:          u.ID,
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
			Status:      GetUserStatus(u.Status),
			Role:        roleCache[u.RoleID],
			CreatedAt:   timestamppb.New(u.CreatedAt),
		}
		users = append(users, user)
	}

	return users, count, hasNextPage, nil
}

func (l *AuthLogic) GetUserByID(ctx context.Context, id uint64) (*models.UserEntity, error) {

	userFromCache := l.RedisStore.GetUser(id)
	if userFromCache != nil {
		return userFromCache, nil
	}

	fmt.Println("reading user from database since cache failed")
	user, err := l.dbStore.AuthRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error in authLogic.getUserByID: %w", err)
	}

	go l.RedisStore.CacheUser(user)

	return user, nil
}
