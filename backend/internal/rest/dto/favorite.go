package dto

import (
	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type FavoriteRequest struct {
	ProfileID uuid.UUID `json:"profile_id" validate:"required,uuid4"`
	QuizID    uuid.UUID `json:"quiz_id" validate:"required,uuid4"`
}

type FavoritesResponse struct {
	Favorites []QuizResponse `json:"favorites"`
}

func ToFavoritesResponse(quizzes []domain.Quiz) *FavoritesResponse {
	favorites := make([]QuizResponse, 0, len(quizzes))
	for _, quiz := range quizzes {
		favorites = append(favorites, *ToQuizResponse(&quiz))
	}

	return &FavoritesResponse{Favorites: favorites}
}
