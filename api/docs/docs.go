// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/header/{block}": {
            "get": {
                "description": "Get ETH Header by block number or hash",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get ETH Header by block",
                "operationId": "get-header-by-block",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Eth header number or hash",
                        "name": "block",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Header"
                        },
                        "headers": {
                            "Token": {
                                "type": "string",
                                "description": "qwerty"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HTTPError"
                        }
                    }
                }
            }
        },
        "/proof/{block}": {
            "get": {
                "description": "Get header with hash proof and mmr roothash",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get header with proof",
                "operationId": "get-header-with-proof",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Eth header number or hash",
                        "name": "block",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.GetEthHeaderWithProofJSONResp"
                        },
                        "headers": {
                            "Token": {
                                "type": "string",
                                "description": "qwerty"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HTTPError"
                        }
                    }
                }
            }
        },
        "/proposal": {
            "post": {
                "description": "Get headers by block numbers, used for relay proposal",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get headers by block numbers",
                "operationId": "get-headers-by-proposal",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "multi",
                        "description": "Eth header numbers",
                        "name": "numbers",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/core.GetEthHeaderWithProofJSONResp"
                                }
                            }
                        },
                        "headers": {
                            "Token": {
                                "type": "string",
                                "description": "qwerty"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HTTPError"
                        }
                    }
                }
            }
        },
        "/receipt/{tx}": {
            "get": {
                "description": "Get receipt by tx hash, used for cross-chain transfer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get receipt by tx hash",
                "operationId": "get-receipt-by-tx",
                "parameters": [
                    {
                        "type": "string",
                        "description": "tx hash",
                        "name": "tx",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "last confirm block",
                        "name": "last",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/core.GetReceiptResp"
                            }
                        },
                        "headers": {
                            "Token": {
                                "type": "string",
                                "description": "qwerty"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "core.GetEthHeaderWithProofJSONResp": {
            "type": "object",
            "properties": {
                "eth_header": {
                    "type": "object",
                    "$ref": "#/definitions/eth.DarwiniaEthHeaderHexFormat"
                },
                "ethash_proof": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/eth.DoubleNodeWithMerkleProof"
                    }
                },
                "mmr_root": {
                    "type": "string"
                }
            }
        },
        "core.GetReceiptResp": {
            "type": "object",
            "properties": {
                "header": {
                    "type": "object",
                    "$ref": "#/definitions/eth.DarwiniaEthHeader"
                },
                "mmr_proof": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "receipt_proof": {
                    "type": "object",
                    "$ref": "#/definitions/eth.ProofRecord"
                }
            }
        },
        "eth.DarwiniaEthHeader": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "difficulty": {
                    "type": "integer"
                },
                "extra_data": {
                    "type": "string"
                },
                "gas_limit": {
                    "type": "integer"
                },
                "gas_used": {
                    "type": "integer"
                },
                "hash": {
                    "type": "string"
                },
                "log_bloom": {
                    "type": "string"
                },
                "number": {
                    "type": "integer"
                },
                "parent_hash": {
                    "type": "string"
                },
                "receipts_root": {
                    "type": "string"
                },
                "seal": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "state_root": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "transactions_root": {
                    "type": "string"
                },
                "uncles_hash": {
                    "type": "string"
                }
            }
        },
        "eth.DarwiniaEthHeaderHexFormat": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "difficulty": {
                    "type": "string"
                },
                "extra_data": {
                    "type": "string"
                },
                "gas_limit": {
                    "type": "string"
                },
                "gas_used": {
                    "type": "string"
                },
                "hash": {
                    "type": "string"
                },
                "log_bloom": {
                    "type": "string"
                },
                "number": {
                    "type": "string"
                },
                "parent_hash": {
                    "type": "string"
                },
                "receipts_root": {
                    "type": "string"
                },
                "seal": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "state_root": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "transactions_root": {
                    "type": "string"
                },
                "uncles_hash": {
                    "type": "string"
                }
            }
        },
        "eth.DoubleNodeWithMerkleProof": {
            "type": "object",
            "properties": {
                "dag_nodes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "proof": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "eth.ProofRecord": {
            "type": "object",
            "properties": {
                "header_hash": {
                    "type": "string"
                },
                "index": {
                    "type": "string"
                },
                "proof": {
                    "type": "string"
                }
            }
        },
        "types.Header": {
            "type": "object",
            "properties": {
                "difficulty": {
                    "type": "string"
                },
                "extraData": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "gasLimit": {
                    "type": "integer"
                },
                "gasUsed": {
                    "type": "integer"
                },
                "logsBloom": {
                    "type": "Bloom"
                },
                "miner": {
                    "type": "string"
                },
                "mixHash": {
                    "type": "string"
                },
                "nonce": {
                    "type": "BlockNonce"
                },
                "number": {
                    "type": "string"
                },
                "parentHash": {
                    "type": "string"
                },
                "receiptsRoot": {
                    "type": "string"
                },
                "sha3Uncles": {
                    "type": "string"
                },
                "stateRoot": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "transactionsRoot": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Swagger Example API",
	Description: "This is a sample server celler server.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
