package dao

// func TestDbmanager(t *testing.T) {
// 	manager, _ := NewDatabase()

// 	patient := auth.Patient{}
// 	patient.PatientAge = 82
// 	patient.PatientName = "谷政理"
// 	patient.PatientPhone = "18118739882"
// 	patient.TreatmentType = "笨"
// 	patient.PatientSex = "女"

// 	patient.PatientUrgentPhone = "18118739880"

// 	if err := manager.Create(&patient); err != nil {
// 		t.Fatal(err)
// 	}

// 	aiderPatient := aider.AiderPatient{}
// 	aiderPatient.Thigh = 0.6
// 	aiderPatient.Calf = 0.6
// 	aiderPatient.Waist = 0.6
// 	aiderPatient.WaistThickness = 0.6
// 	aiderPatient.Diagnosis = "没救"
// 	aiderPatient.PatientID = patient.ID

// 	if err := manager.Create(&aiderPatient); err != nil {
// 		t.Fatal(err)
// 	}

// 	user := auth.User{}
// 	user.UserID = "lmq"
// 	user.UserName = "梁敏弱"
// 	user.Comments = "测试"
// 	user.Email = "liangminqiang@buffalo-robot.com"
// 	user.Address = "上海"
// 	user.HashPassword("123456")
// 	user.Phone = "18111111111"

// 	if err := manager.Create(&user); err != nil {
// 		t.Fatal(err)
// 	}

// }

// func TestMigrage(t *testing.T) {
// 	manager, _ := NewDatabase()

// 	manager.Migrate()
// }

// func TestAddUserPatientBind(t *testing.T) {
// 	manager, _ := NewDatabase()

// 	var userbind []auth.UserPatientBind
// 	if err := manager.DB.Find(&userbind).Error; err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println(userbind[0])
// 	}
// }

// func TestDeviceBindInfo(t *testing.T) {
// 	manager, _ := NewDatabase()

// 	bindinfo := models.DeviceUserBind{}
// 	bindinfo.Comments = "绑定测试"
// 	bindinfo.UserModelID = "Admin"
// 	bindinfo.DeviceModelID = "5ebba18db36ccdf8cc81e470"

// 	if err := manager.Create(&bindinfo); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestMongoUpdate(t *testing.T) {
// 	manager, _ := NewDatabase()

// 	type UpdateTest struct {
// 		ID    primitive.ObjectID `bson:"_id,omitempty"`
// 		One   string             `bson:"one"`
// 		Two   float64            `bson:"two"`
// 		Three int64              `bson:"three"`
// 	}

// 	first := UpdateTest{
// 		One:   "haha",
// 		Two:   22.4,
// 		Three: 2,
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	query, err := manager.GetMongo().Database(constants.MongoDB).Collection("test").InsertOne(ctx, &first)
// 	first.ID = query.InsertedID.(primitive.ObjectID)
// 	fmt.Println(first.ID.Hex())

// 	updateOne := UpdateTest{
// 		One:   "haha5",
// 		Two:   22.44,
// 		Three: 52,
// 	}

// 	ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel3()

// 	var read UpdateTest
// 	err = manager.GetMongo().Database(constants.MongoDB).Collection("test").FindOne(ctx3, bson.M{"_id": first.ID}).Decode(&read)
// 	fmt.Println(err)
// 	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel2()
// 	_, err = manager.GetMongo().Database(constants.MongoDB).Collection("test").ReplaceOne(ctx2, bson.M{"_id": first.ID}, updateOne)
// 	fmt.Println(err)
// }

// func TestAuthCreateTable(t *testing.T) {
// 	CreateRoleTable()
// }
