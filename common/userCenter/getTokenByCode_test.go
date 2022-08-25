package userCenter

import "testing"


func TestGetTokenByCode(t *testing.T) {
	info, err := userC.GetTokenByCode(&GetTokenByCode{
		Code:         "8d37fcdf52b2355e8c0513afea6c6085",
		//Code:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZU51bWJlciI6IjE5OTExMTEyMjIyIiwiaWQiOiJkY2RlZDBhNmRlMTg0MDhiOTAyYzNhY2JhZWJjODEyZSIsImlhdCI6MTY1NDYxNTMwNywiY2xpZW50SWQiOiJQai1EaWFnbm9zaXMtQyJ9.8msFEhYxaLNf5Q-GYhVIYMZU2avh-1O4SriOTjjNozM",
		Grant_type:   "authorization_code",
		Three_type:   "1",
	})
	if err != nil {
		panic(err)
	}
	t.Logf("%+v, %+v",info, err)
}

//{
//	Code:200
//	Msg:successful
//	Data:{
//		CliendId:Pj-Diagnosis-C
//		AccessToken:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2MjY5OTcsInN1YiI6ImRjZGVkMGE2ZGUxODQwOGI5MDJjM2FjYmFlYmM4MTJlIiwiaWF0IjoxNjUyMzY3Nzk3LCJ0eXBlIjoiYWNjZXNzIiwianRpIjoidXNlciIsImNsIjoiUGotRGlhZ25vc2lzLUMifQ.XsEKQaKYFaZrOTQoRFdNvWOHehxD82n9LadHELQrNZg RefreshToken:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI5NzI1OTcsInN1YiI6ImRjZGVkMGE2ZGUxODQwOGI5MDJjM2FjYmFlYmM4MTJlIiwiaWF0IjoxNjUyMzY3Nzk3LCJ0eXBlIjoicmVmcmVzaCIsImp0aSI6InVzZXIiLCJjbCI6IlBqLURpYWdub3Npcy1DIn0.Dq_rEYu1kXMnIe-Od8eIw8f35K3wCzg3Ofr8FG0Sca0
//		AccessExpireIn:259200
//		RefreshExpireIn:604800
//		TokenType:bearer
//		Info:{
//			CreateTime:2022-05-12 14:33:55
//			UpdateTime:2022-05-12 14:33:55
//			Id:dcded0a6de18408b902c3acbaebc812e
//			RealName:sss
//			Gender:男
//			Password:$2b$12$mhCaNK36vtWCmmH1EZwWU.8/PtmXalqGhFoiWRS0GK6pwaadgL31G
//			Superuser:0
//			PhoneNumber:19911112222
//			Flag:0
//			UserInfo:{
//				CreateTime:2022-05-12 14:33:55
//				UpdateTime:2022-05-12 14:33:55
//				Id:133
//				UserId:dcded0a6de18408b902c3acbaebc812e
//				Server:Pj-Diagnosis-C
//				NickName:1
//				UserEmail:111@163.com
//				OtherContact:
//				ExpertImage:
//				Hospital:sss
//				Department:sddsvs
//				JobNumber:54534
//				OwnerUser:
//			}
//		}
//	}
//}

//{Code:200 Msg:已获取授权 Data:{CliendId:Pj-Diagnosis-C AccessToken:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2MjY5OTcsInN1YiI6ImRjZGVkMGE2ZGUxODQwOGI5MDJjM2FjYmFlYmM4MTJlIiwiaWF0IjoxNjUyMzY3Nzk3LCJ0eXBlIjoiYWNjZXNzIiwianRpIjoidXNlciIsImNsIjoiUGotRGlhZ25vc2lzLUMifQ.XsEKQaKYFaZrOTQoRFdNvWOHehxD82n9LadHELQrNZg RefreshToken:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI5NzI1OTcsInN1YiI6ImRjZGVkMGE2ZGUxODQwOGI5MDJjM2FjYmFlYmM4MTJlIiwiaWF0IjoxNjUyMzY3Nzk3LCJ0eXBlIjoicmVmcmVzaCIsImp0aSI6InVzZXIiLCJjbCI6IlBqLURpYWdub3Npcy1DIn0.Dq_rEYu1kXMnIe-Od8eIw8f35K3wCzg3Ofr8FG0Sca0 AccessExpireIn:259200 RefreshExpireIn:604800 TokenType:bearer Info:{CreateTime:2022-05-12 14:33:55 UpdateTime:2022-05-12 14:33:55 Id:dcded0a6de18408b902c3acbaebc812e RealName:sss Gender:男 Password:$2b$12$mhCaNK36vtWCmmH1EZwWU.8/PtmXalqGhFoiWRS0GK6pwaadgL31G Superuser:0 PhoneNumber:19911112222 Flag:0 UserInfo:{CreateTime:2022-05-12 14:33:55 UpdateTime:2022-05-12 14:33:55 Id:133 UserId:dcded0a6de18408b902c3acbaebc812e Server:Pj-Diagnosis-C NickName: UserEmail:111@163.com OtherContact: ExpertImage: Hospital:sss Department:sddsvs JobNumber:54534 OwnerUser:}}}}