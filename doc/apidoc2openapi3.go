package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type APIInfo struct {
	APIURL           string      `json:"api_url"`
	HTTPMethod       string      `json:"http_method"`
	ShortDescription string      `json:"short_description"`
	Deprecated       interface{} `json:"deprecated"`
}

type Param struct {
	Name         string        `json:"name"`
	FullName     string        `json:"full_name"`
	Description  string        `json:"description"`
	Required     bool          `json:"required"`
	AllowNil     bool          `json:"allow_nil"`
	Validator    string        `json:"validator"`
	ExpectedType string        `json:"expected_type"`
	Metadata     interface{}   `json:"metadata"`
	Show         bool          `json:"show"`
	Validations  []interface{} `json:"validations"`
	Params       []Param       `json:"params,omitempty"`
}

type Example struct {
	Verb         string      `json:"verb"`
	Path         string      `json:"path"`
	Versions     []string    `json:"versions"`
	Query        string      `json:"query"`
	RequestData  interface{} `json:"request_data"`
	ResponseData interface{} `json:"response_data"`
	Code         string      `json:"code"`
	ShowInDoc    int         `json:"show_in_doc"`
	Recorded     bool        `json:"recorded"`
}

type MetadataRoles struct {
	RequiredScopes []string `json:"required_scopes"`
	Roles          []string `json:"roles"`
}

type Method struct {
	DocURL          string        `json:"doc_url"`
	Name            string        `json:"name"`
	Apis            []APIInfo     `json:"apis"`
	Formats         interface{}   `json:"formats"`
	FullDescription string        `json:"full_description"`
	Errors          []interface{} `json:"errors"`
	Params          []Param       `json:"params"`
	Examples        []Example     `json:"examples"`
	Metadata        MetadataRoles `json:"metadata,omitempty"`
	See             []interface{} `json:"see"`
	Headers         []interface{} `json:"headers"`
	Show            bool          `json:"show"`
}

type Resource struct {
	DocURL           string        `json:"doc_url"`
	APIURL           string        `json:"api_url"`
	Name             string        `json:"name"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	Version          string        `json:"version"`
	Formats          interface{}   `json:"formats"`
	Metadata         interface{}   `json:"metadata"`
	Methods          []Method      `json:"methods"`
	Headers          []interface{} `json:"headers"`
}

type Resources map[string]Resource

type ApiDoc struct {
	Docs struct {
		Name      string    `json:"name"`
		Info      string    `json:"info"`
		Copyright string    `json:"copyright"`
		DocURL    string    `json:"doc_url"`
		APIURL    string    `json:"api_url"`
		Resources Resources `json:"resources"`
	} `json:"docs"`
}

func initializeOpenAPI(apiDoc *ApiDoc, commonPrefix string) map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       apiDoc.Docs.Name,
			"description": apiDoc.Docs.Info,
			"version":     "1.0.0",
		},
		"servers": []map[string]interface{}{
			{
				"url": commonPrefix,
			},
		},
		"paths":    map[string]interface{}{},
		"consumes": []string{"application/json"},
		"produces": []string{"application/json"},
	}
}

func findCommonPrefix(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	// Find the shortest path to use as a starting point
	shortestPath := paths[0]
	for _, path := range paths {
		if len(path) < len(shortestPath) {
			shortestPath = path
		}
	}

	// Iterate over the shortest path to find the common prefix
	for i, c := range shortestPath {
		for _, path := range paths {
			if path[i] != byte(c) {
				commonPrefix := shortestPath[:i]
				if strings.HasSuffix(commonPrefix, "/") {
					commonPrefix = commonPrefix[:len(commonPrefix)-1]
				}
				return commonPrefix
			}
		}
	}

	return shortestPath
}

func generateOperationID(httpMethod, path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	var operationID []string
	operationID = append(operationID, strings.Title(httpMethod))

	for i, segment := range segments {
		// Add "From" or "By" before path parameters, e.g., :achievement_id
		if strings.HasPrefix(segment, ":") {
			segment = segment[1:]
			if i < len(segments)-1 {
				operationID = append(operationID, "From")
			} else {
				operationID = append(operationID, "By")
			}
		}

		words := regexp.MustCompile(`([A-Za-z0-9]+)`).FindAllString(segment, -1)

		for _, word := range words {
			operationID = append(operationID, strings.Title(strings.ToLower(word)))
		}
	}

	return strings.Join(operationID, "")
}

func convertApiDocToOpenAPI(apiDoc *ApiDoc) map[string]interface{} {
	// Find the common prefix of all the paths
	var paths []string
	for _, resource := range apiDoc.Docs.Resources {
		for _, method := range resource.Methods {
			for _, api := range method.Apis {
				paths = append(paths, api.APIURL)
			}
		}
	}
	commonPrefix := findCommonPrefix(paths)
	openAPI := initializeOpenAPI(apiDoc, commonPrefix)

	// Define the Response schema
	responseSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type": "string",
			},
			"message": map[string]interface{}{
				"type": "string",
			},
		},
	}

	// Add the Response schema to the components section
	components := map[string]interface{}{
		"schemas": map[string]interface{}{
			"Response": responseSchema,
		},
	}
	openAPI["components"] = components

	// Add the paths to the openAPI
	for _, resource := range apiDoc.Docs.Resources {
		for _, method := range resource.Methods {
			addMethodToOpenAPIPaths(openAPI, commonPrefix, resource, method)
		}
	}

	return openAPI
}

func addMethodToOpenAPIPaths(openAPI map[string]interface{}, commonPrefix string, resource Resource, method Method) {
	paths := openAPI["paths"].(map[string]interface{})

	for _, api := range method.Apis {
		path := strings.TrimPrefix(api.APIURL, commonPrefix)
		httpMethod := strings.ToLower(api.HTTPMethod)

		if _, exists := paths[path]; !exists {
			paths[path] = map[string]interface{}{}
		}

		operation := map[string]interface{}{
			"tags":        []string{resource.Name},
			"summary":     method.Name,
			"description": method.FullDescription,
			"operationId": generateOperationID(httpMethod, path),
			"parameters":  []map[string]interface{}{},
			"responses":   map[string]interface{}{},
		}

		if api.Deprecated != nil {
			operation["deprecated"] = true
		}

		// Add parameters
		for _, param := range method.Params {
			parameters := operation["parameters"].([]map[string]interface{})
			swaggerType := convertToOpenAPIType(param.ExpectedType)
			parameter := map[string]interface{}{
				"name": param.Name,
				"in":   "query",
				"type": swaggerType,
			}
			if swaggerType == "array" {
				parameter["items"] = map[string]interface{}{
					"type": "string",
				}
				parameter["collectionFormat"] = "csv"
			}
			parameters = append(parameters, parameter)
			operation["parameters"] = parameters
		}

		// Add responses
		for _, example := range method.Examples {
			responseCode := example.Code

			num, err := strconv.Atoi(responseCode)
			if err != nil {
				continue
			}

			// Check if the response code is a valid HTTP status code
			if http.StatusText(num) == "" {
				continue
			}

			responses := operation["responses"].(map[string]interface{})
			response := map[string]interface{}{
				"description": "",
			}

			if example.ResponseData != nil {
				response["content"] = map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema":  generateSchemaFromExample(example.ResponseData),
						"example": example.ResponseData,
					},
				}
			}

			responses[responseCode] = response
		}

		// Add a default response if there are no valid responses
		if len(operation["responses"].(map[string]interface{})) == 0 {
			operation["responses"] = map[string]interface{}{
				"default": map[string]interface{}{
					"description": "Default response",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/Response",
							},
						},
					},
				},
			}
		}

		// Add the operation to the path
		pathItem := paths[path].(map[string]interface{})
		pathItem[httpMethod] = operation
	}
}

func generateSchemaFromExample(example interface{}) map[string]interface{} {
	if example == nil {
		return map[string]interface{}{
			"type": "null",
		}
	}

	switch exampleType := example.(type) {
	case []interface{}:
		if len(exampleType) > 0 {
			return map[string]interface{}{
				"type":  "array",
				"items": generateSchemaFromExample(exampleType[0]),
			}
		}
		return map[string]interface{}{
			"type": "array",
		}
	case map[string]interface{}:
		properties := make(map[string]interface{})
		required := []string{}
		for key, value := range exampleType {
			properties[key] = generateSchemaFromExample(value)
			required = append(required, key)
		}
		return map[string]interface{}{
			"type":       "object",
			"properties": properties,
			"required":   required,
		}
	case float64:
		return map[string]interface{}{
			"type": "number",
		}
	case string:
		return map[string]interface{}{
			"type": "string",
		}
	case bool:
		return map[string]interface{}{
			"type": "boolean",
		}
	default:
		return map[string]interface{}{}
	}
}

func convertToOpenAPIType(apiDocType string) string {
	switch apiDocType {
	case "integer":
		return "integer"
	case "float":
		return "number"
	case "boolean":
		return "boolean"
	case "hash":
		return "string"
	case "array":
		return "array"
	default:
		return "string"
	}
}

func parseApiDocFromFile(filename string) (*ApiDoc, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var apiDoc ApiDoc
	err = json.Unmarshal(data, &apiDoc)
	if err != nil {
		return nil, err
	}

	return &apiDoc, nil
}

func saveOpenAPIFile(filename string, openAPI map[string]interface{}) error {
	data, err := json.MarshalIndent(openAPI, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var apidocFile string
	var swaggerFile string

	// Define command-line flags
	flag.StringVar(&apidocFile, "apidoc", "apidoc.json", "Path to the input apidoc file")
	flag.StringVar(&swaggerFile, "swagger", "swagger.json", "Path to the output swagger file")

	// Parse command-line flags
	flag.Parse()

	// Show usage information if no arguments are provided
	if apidocFile == "" || swaggerFile == "" {
		flag.Usage()
		return
	}

	// Process files
	apiDoc, err := parseApiDocFromFile(apidocFile)
	if err != nil {
		log.Fatalf("Error parsing apidoc file: %s", err)
	}

	if err = saveOpenAPIFile(swaggerFile, convertApiDocToOpenAPI(apiDoc)); err != nil {
		log.Fatalf("Error generating swagger file: %s", err)
	}

	log.Printf("Swagger JSON generated: %s\n", swaggerFile)
}
