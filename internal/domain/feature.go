package domain

type SummarizeRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	Mode       string `json:"mode" binding:"required,oneof=short detailed mindmap"`
}

type QuizRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=multiple_choice essay true_false"`
	Count      int    `json:"count" binding:"required,min=1,max=20"`
	Difficulty string `json:"difficulty" binding:"required,oneof=easy medium hard"`
}

type EssayRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required,min=50"`
}
