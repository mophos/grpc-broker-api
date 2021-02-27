package main

import (
	"context"
	"fmt"
	"log"
	"os"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/mophos/grpc-broker-api/database"

	// "github.com/mophos/grpc-broker-api/proto"
	proto "github.com/moph-gateway/his-proto/proto"
	"github.com/mophos/grpc-broker-api/user"
	"google.golang.org/grpc"
)

func initDatabase() {
	var err error
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USERNAME")
	db_pass := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	url := db_user + ":" + db_pass + "@tcp(" + db_host + ":" + db_port + ")/" + db_name + "?charset=utf8&parseTime=True"
	database.DBConn, err = gorm.Open("mysql", url)
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
}

func main() {
	err := godotenv.Load("conf.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	initDatabase()
	defer database.DBConn.Close()

	//hosxpv3
	urlHosxpv3 := os.Getenv("URL_HOSXPV3")
	urlHosxpv4 := os.Getenv("URL_HOSXPV4")
	urlHosxppcu := os.Getenv("URL_HOSXPPCU")
	fmt.Print(urlHosxpv3, urlHosxpv4, urlHosxppcu)
	connHosxpv3, err := grpc.Dial(urlHosxpv3, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	clientEmrHosxpv3 := proto.NewEmrServiceClient(connHosxpv3)
	clientMasterHosxpv3 := proto.NewMasterServiceClient(connHosxpv3)

	//hosxpv4
	connHosxpv4, err := grpc.Dial(urlHosxpv4, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	clientEmrHosxpv4 := proto.NewEmrServiceClient(connHosxpv4)
	clientMasterHosxpv4 := proto.NewMasterServiceClient(connHosxpv4)

	//hosxppcu
	// connHosxppcu, err := grpc.Dial(urlHosxppcu, grpc.WithInsecure())
	// if err != nil {
	// 	panic(err)
	// }
	// clientEmrHosxppcu := proto.NewEmrServiceClient(connHosxpv4)
	// clientMasterHosxppcu := proto.NewMasterServiceClient(connHosxppcu)

	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	api := app.Group("api")
	api.Post("/v1/login", user.Login)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	api.Get("/jwt", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		return c.SendString("Welcome " + name)
	})
	api.Post("/v1/patient-info", func(c *fiber.Ctx) error {
		cid := c.FormValue("cid")

		req := &proto.RequestCid{Cid: cid}

		if res, err := clientEmrHosxpv3.PatientInfo(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"results": res.Results,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	})

	api.Post("/v1/services", func(c *fiber.Ctx) error {
		cid := c.FormValue("cid")

		req := &proto.RequestCid{Cid: cid}

		if res, err := clientEmrHosxpv3.GetServices(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"results": res.Results,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	})

	api.Post("/v1/doctor", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")

		req := &proto.RequestHospcode{Hospcode: hospcode}

		res, err := clientMasterHosxpv3.DoctorList(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientMasterHosxpv4.DoctorList(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientMasterHosxppcu.DoctorList(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})
	api.Post("/v1/clinic", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")

		req := &proto.RequestHospcode{Hospcode: hospcode}

		res, err := clientMasterHosxpv3.ClinicList(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientMasterHosxpv4.ClinicList(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientMasterHosxppcu.ClinicList(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/v1/screening", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")
		hn := c.FormValue("hn")
		vn := c.FormValue("vn")

		req := &proto.RequestPatient{Hospcode: hospcode, Hn: hn, Vn: vn}

		res, err := clientEmrHosxpv3.GetScreening(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientEmrHosxpv4.GetScreening(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientEmrHosxppcu.GetScreening(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/v1/diagnosis", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")
		hn := c.FormValue("hn")
		vn := c.FormValue("vn")

		req := &proto.RequestPatient{Hospcode: hospcode, Hn: hn, Vn: vn}

		res, err := clientEmrHosxpv3.GetDiagnosis(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientEmrHosxpv4.GetDiagnosis(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientEmrHosxppcu.GetDiagnosis(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/v1/procedure", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")
		hn := c.FormValue("hn")
		vn := c.FormValue("vn")

		req := &proto.RequestPatient{Hospcode: hospcode, Hn: hn, Vn: vn}

		res, err := clientEmrHosxpv3.GetProcedure(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientEmrHosxpv4.GetProcedure(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientEmrHosxppcu.GetProcedure(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/v1/lab", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")
		hn := c.FormValue("hn")
		vn := c.FormValue("vn")

		req := &proto.RequestPatient{Hospcode: hospcode, Hn: hn, Vn: vn}

		res, err := clientEmrHosxpv3.GetLab(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientEmrHosxpv4.GetLab(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientEmrHosxppcu.GetLab(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/v1/vaccine", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")
		hn := c.FormValue("hn")
		vn := c.FormValue("vn")

		req := &proto.RequestPatient{Hospcode: hospcode, Hn: hn, Vn: vn}

		res, err := clientEmrHosxpv3.GetVaccine(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientEmrHosxpv4.GetVaccine(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientEmrHosxppcu.GetVaccine(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/v1/drug", func(c *fiber.Ctx) error {
		hospcode := c.FormValue("hospcode")
		hn := c.FormValue("hn")
		vn := c.FormValue("vn")

		req := &proto.RequestPatient{Hospcode: hospcode, Hn: hn, Vn: vn}

		res, err := clientEmrHosxpv3.GetDrug(context.Background(), req)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}

		res2, err2 := clientEmrHosxpv4.GetDrug(context.Background(), req)
		if err2 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		// res3, err3 := clientEmrHosxppcu.GetDrug(context.Background(), req)
		// if err3 != nil {
		// 	log.Fatalf("open stream error %v", err2)
		// }

		data := append(res.Results, res2.Results...)
		// data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf("0.0.0.0:%s", port)))
}
