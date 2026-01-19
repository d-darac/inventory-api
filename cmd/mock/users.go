package main

// func createUser(q *database.Queries) (*database.CreateUserRow, error) {
// 	n := time.Now().UnixNano()
// 	password := fmt.Sprintf("super_strong_password_%d", n)
// 	hashedPass, err := auth.HashPassword(password)
// 	if err != nil {
// 		return nil, fmt.Errorf("couldn't hash password: %v", err)
// 	}
// 	user, err := q.CreateUser(context.Background(), database.CreateUserParams{
// 		Email:          fmt.Sprintf("email_%d@test.com", n),
// 		HashedPassword: hashedPass,
// 		Name: sql.NullString{
// 			String: fmt.Sprintf("Test User %d", n),
// 			Valid:  true,
// 		},
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("couldn't create test user: %v", err)
// 	}
// 	return &user, nil
// }
