package requests

type UpdateProfileParam struct {
	Name            string `json:"name"`
	Email           string `json:"email" binding:"omitempty,email"`
	Password        string `json:"password"`
	ProfileImageUrl string `json:"profile_image_url" binding:"omitempty,url"`
}
