package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/siteslave/grpc-rest-client/proto"
	"google.golang.org/grpc"
)

func main() {
	//hosxpv3
	connHosxpv3, err := grpc.Dial("localhost:4042", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	clientEmrHosxpv3 := proto.NewEmrServiceClient(connHosxpv3)
	clientMasterHosxpv3 := proto.NewMasterServiceClient(connHosxpv3)

	//hosxpv4
	connHosxpv4, err := grpc.Dial("localhost:4043", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// clientEmrHosxpv4 := proto.NewEmrServiceClient(connHosxpv4)
	clientMasterHosxpv4 := proto.NewMasterServiceClient(connHosxpv4)

	//hosxppcu
	connHosxppcu, err := grpc.Dial("localhost:4044", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// clientEmrHosxpv4 := proto.NewEmrServiceClient(connHosxpv4)
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

	log.Fatal(app.Listen(":3003"))
}
