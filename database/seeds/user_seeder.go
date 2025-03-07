package seeds

import (
	"blog/internal/core/domain/model"
	"blog/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB) {
	bytes, err := conv.HashPassword("password")
	if err != nil {
		log.Fatal().Msgf("Error creating password")
	}

	admin := model.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@gmail.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Error seending admin role")
	} else {
		log.Info().Msg("Admin role seeded successfully")
	}
}
