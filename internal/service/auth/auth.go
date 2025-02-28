package auth

import (
	"context"

	"github.com/segmentfault/answer/internal/entity"
	"github.com/segmentfault/answer/pkg/token"
	"github.com/segmentfault/pacman/log"
)

// AuthRepo auth repository
type AuthRepo interface {
	GetUserCacheInfo(ctx context.Context, accessToken string) (userInfo *entity.UserCacheInfo, err error)
	SetUserCacheInfo(ctx context.Context, accessToken string, userInfo *entity.UserCacheInfo) error
	RemoveUserCacheInfo(ctx context.Context, accessToken string) (err error)
	GetUserStatus(ctx context.Context, userID string) (userInfo *entity.UserCacheInfo, err error)
	GetCmsUserCacheInfo(ctx context.Context, accessToken string) (userInfo *entity.UserCacheInfo, err error)
	SetCmsUserCacheInfo(ctx context.Context, accessToken string, userInfo *entity.UserCacheInfo) error
	RemoveCmsUserCacheInfo(ctx context.Context, accessToken string) (err error)
}

// AuthService kit service
type AuthService struct {
	authRepo AuthRepo
}

// NewAuthService email service
func NewAuthService(authRepo AuthRepo) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (as *AuthService) GetUserCacheInfo(ctx context.Context, accessToken string) (userInfo *entity.UserCacheInfo, err error) {
	userCacheInfo, err := as.authRepo.GetUserCacheInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	cacheInfo, _ := as.authRepo.GetUserStatus(ctx, userCacheInfo.UserID)
	if cacheInfo != nil {
		log.Infof("user status updated: %+v", cacheInfo)
		userCacheInfo.UserStatus = cacheInfo.UserStatus
		userCacheInfo.EmailStatus = cacheInfo.EmailStatus
		// update current user cache info
		err := as.authRepo.SetUserCacheInfo(ctx, accessToken, userCacheInfo)
		if err != nil {
			return nil, err
		}
	}
	return userCacheInfo, nil
}

func (as *AuthService) SetUserCacheInfo(ctx context.Context, userInfo *entity.UserCacheInfo) (accessToken string, err error) {
	accessToken = token.GenerateToken()
	err = as.authRepo.SetUserCacheInfo(ctx, accessToken, userInfo)
	return accessToken, err
}

func (as *AuthService) RemoveUserCacheInfo(ctx context.Context, accessToken string) (err error) {
	return as.authRepo.RemoveUserCacheInfo(ctx, accessToken)
}

//cms

func (as *AuthService) GetCmsUserCacheInfo(ctx context.Context, accessToken string) (userInfo *entity.UserCacheInfo, err error) {
	return as.authRepo.GetCmsUserCacheInfo(ctx, accessToken)
}

func (as *AuthService) SetCmsUserCacheInfo(ctx context.Context, accessToken string, userInfo *entity.UserCacheInfo) (err error) {
	err = as.authRepo.SetCmsUserCacheInfo(ctx, accessToken, userInfo)
	return err
}

func (as *AuthService) RemoveCmsUserCacheInfo(ctx context.Context, accessToken string) (err error) {
	return as.authRepo.RemoveCmsUserCacheInfo(ctx, accessToken)
}
