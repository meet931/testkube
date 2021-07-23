package v1

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/kubeshop/kubetest/pkg/api/kubetest"
	"github.com/kubeshop/kubetest/pkg/executor/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) GetAllScripts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("OK 👋!")
	}
}

func (s Server) ExecuteScript() fiber.Handler {
	// TODO use kube API to get registered executor details - for now it'll be fixed
	// we need to choose client based on script type in future for now there is only
	// one client postman-collection newman based executor
	// should be done on top level from some kind of available clients poll
	// consider moving them to separate struct - and allow to choose by executor ID
	executorClient := client.NewHTTPExecutorClient(client.DefaultURI)

	return func(c *fiber.Ctx) error {

		scriptID := c.Params("id")
		s.Log.Infow("running execution of script", "id", scriptID)

		var request struct {
			Name string
		}
		c.BodyParser(&request)

		// TODO use kubeapi to get script content
		content := exampleCollection
		execution, err := executorClient.Execute(content)
		if err != nil {
			return err
		}

		ctx := context.Background()
		scriptExecution := kubetest.NewScriptExecution(
			primitive.NewObjectID().Hex(),
			request.Name,
			execution,
		)
		s.Repository.Insert(ctx, scriptExecution)

		execution, err = executorClient.Watch(execution.Id, func(e kubetest.Execution) error {
			s.Log.Infow("saving", "status", e.Status, "execution", e)
			scriptExecution.Execution = &e
			return s.Repository.Update(ctx, scriptExecution)
		})

		if err != nil {
			return err
		}

		return c.JSON(execution)
	}
}

func (s Server) GetAllScriptExecutions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("OK 👋!")
	}
}

func (s Server) GetScriptExecution() fiber.Handler {
	// TODO use kube API to get registered executor details - for now it'll be fixed
	// we need to choose client based on script type in future for now there is only
	// one client postman-collection newman based executor
	// should be done on top level from some kind of available clients poll
	// consider moving them to separate struct - and allow to choose by executor ID

	executorClient := client.NewHTTPExecutorClient(client.DefaultURI)
	return func(c *fiber.Ctx) error {

		scriptID := c.Params("id")
		executionID := c.Params("executionID")
		s.Log.Infow("GET execution request", "id", scriptID, "executionID", executionID)

		execution, err := executorClient.Get(executionID)
		if err != nil {
			return err
		}
		return c.JSON(execution)
	}
}

func (s Server) AbortScriptExecution() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

// TODO remove when reading from API will be implemented
const exampleCollection = `{
	"info": {
		"_postman_id": "fa1ce97f-ff5d-40ed-9c9c-e0a92063ce98",
		"name": "Remotes",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	},
	"item": [
		{
			"name": "Google",
			"event": [
				{
					"listen": "test",
					"script": {
						"exec": [
							"    pm.test(\"Successful GET request\", function () {",
							"        pm.expect(pm.response.code).to.be.oneOf([200, 201, 202]);",
							"    });"
						],
						"type": "text/javascript"
					}
				}
			],
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "https://google.com",
					"protocol": "https",
					"host": [
						"google",
						"com"
					]
				}
			},
			"response": []
		},
		{
			"name": "Kasia.in Homepage",
			"event": [
				{
					"listen": "test",
					"script": {
						"exec": [
							"pm.test(\"Body matches string\", function () {",
							"    pm.expect(pm.response.text()).to.include(\"PRZEPIS NA CHLEB\");",
							"});"
						],
						"type": "text/javascript"
					}
				}
			],
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "https://kasia.in",
					"protocol": "https",
					"host": [
						"kasia",
						"in"
					]
				}
			},
			"response": []
		}
	]
}`