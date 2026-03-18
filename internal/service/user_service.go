package service

import (
	"errors"

	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
func (s *UserService) GetMe(userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return user, nil
}

func (s *UserService) UpdateProfile(userID string, req *domain.EditProfileRequest) (*domain.User, error) {
	// 1. Cari user berdasarkan ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// 2. Jika email ingin diubah, cek apakah email baru sudah terpakai
	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != "" {
			return nil, errors.New("email sudah digunakan oleh pengguna lain")
		}
		user.Email = req.Email
	}

	// 3. Update nama jika diisi
	if req.Name != "" {
		user.Name = req.Name
	}

	// 4. Simpan perubahan ke database
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, errors.New("gagal menyimpan pembaruan profil")
	}

	return user, nil
}

// ChangePassword menangani logika pergantian password
func (s *UserService) ChangePassword(userID string, req *domain.ChangePasswordRequest) error {
	// 1. Cari user berdasarkan ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// 2. Verifikasi password lama
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errors.New("password lama salah")
	}

	// 3. Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	// 4. Update dan simpan password baru ke database
	user.Password = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return errors.New("gagal menyimpan password baru")
	}

	return nil
}
