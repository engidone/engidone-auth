package greet

type greetingRepository struct{}

func NewGreetingRepository() *greetingRepository {
	return &greetingRepository{}
}
