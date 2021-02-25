package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/siteslave/grpc-rest-client/proto"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

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
	connHosxppcu, err := grpc.Dial(urlHosxppcu, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	clientEmrHosxppcu := proto.NewEmrServiceClient(connHosxpv4)
	clientMasterHosxppcu := proto.NewMasterServiceClient(connHosxppcu)

	app := fiber.New()
	app.Use(logger.New())
	app.Post("/patient-info", func(c *fiber.Ctx) error {
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

	app.Post("/services", func(c *fiber.Ctx) error {
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

	app.Post("/doctor", func(c *fiber.Ctx) error {
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

		res3, err3 := clientMasterHosxppcu.DoctorList(context.Background(), req)
		if err3 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		data := append(res.Results, res2.Results...)
		data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})
	app.Post("/clinic", func(c *fiber.Ctx) error {
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

		res3, err3 := clientMasterHosxppcu.ClinicList(context.Background(), req)
		if err3 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		data := append(res.Results, res2.Results...)
		data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	app.Post("/screening", func(c *fiber.Ctx) error {
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

		res3, err3 := clientEmrHosxppcu.GetScreening(context.Background(), req)
		if err3 != nil {
			log.Fatalf("open stream error %v", err2)
		}

		data := append(res.Results, res2.Results...)
		data = append(data, res3.Results...)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"results": data,
		})

		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })

	})

	log.Fatal(app.Listen(":3003"))
}
