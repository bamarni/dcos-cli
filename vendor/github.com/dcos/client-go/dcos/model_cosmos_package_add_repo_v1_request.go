/*
 * DC/OS
 *
 * DC/OS API
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dcos

type CosmosPackageAddRepoV1Request struct {
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Index int32  `json:"index,omitempty"`
}
