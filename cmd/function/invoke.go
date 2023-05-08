package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/morty-faas/cli/cliconfig"
	"github.com/morty-faas/cli/pkg/debug"
	"github.com/morty-faas/cli/pkg/httpclient"

	"github.com/oliveagle/jsonpath"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type invokeOptions struct {
	FnName    string   `json:"functionName"`
	FnVersion string   `json:"functionVersion"`
	Method    string   `json:"method"`
	Body      string   `json:"body"`
	Headers   []string `json:"headers"`
	Params    []string `json:"params"`
}

const (
	invokeFunctionEndpoint = "functions/{name}/{version}/invoke"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke NAME",
	Short: "Invoke a function",
	Long:  `Invoke a function using default options or choose HTTP method, body, headers etc.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdCtx := cmd.Context()

		// Safe call, validation is performed by cobra.ExactArgs(1) above
		name := args[0]
		version, _ := cmd.Flags().GetString("version")

		// Extract the method from the command
		method, _ := cmd.Flags().GetString("method")
		method = strings.ToUpper(method)
		if err := validateHttpMethod(method); err != nil {
			return err
		}

		data, _ := cmd.Flags().GetString("data")
		headers, _ := cmd.Flags().GetStringArray("headers")
		jsonPathQuery, _ := cmd.Flags().GetString("query")
		params, _ := cmd.Flags().GetStringArray("param")

		opts := &invokeOptions{
			FnName:    name,
			FnVersion: version,
			Method:    method,
			Body:      data,
			Headers:   headers,
			Params:    params,
		}

		response, err := invoke(cmdCtx, opts)
		if err != nil {
			return err
		}

		// Try to decode the response as a JSON object
		// If an errror occurs, it means that the response
		// is simply a string so we can print it directly
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(response), &jsonData); err != nil {
			fmt.Println(response)
			return nil
		}

		// Try to apply the JSONPath query
		output, err := jsonpath.JsonPathLookup(jsonData, jsonPathQuery)
		if err != nil {
			return err
		}

		// If the value is a JSON, encode it
		if v, ok := output.(map[string]interface{}); ok {
			by, _ := json.Marshal(v)
			output = string(by)
		}

		// We don't want to print an output if the function
		// returns an empty payload
		if output != "null" {
			fmt.Println(output)
		}

		return nil
	},
}

func init() {
	invokeCmd.PersistentFlags().String("version", "latest", "The version of the function to invoke.")
	invokeCmd.PersistentFlags().StringP("method", "X", "GET", "The HTTP method to use to invoke the request. Valid values are: GET, POST, PUT, PATCH, DELETE")
	invokeCmd.PersistentFlags().StringP("data", "d", "", "The body to pass to the invocation request.")
	invokeCmd.PersistentFlags().StringP("query", "q", "$", "A valid JSON Path expression to execute on the function response. If the function output isn't a JSON, the flag will have no effect.")
	invokeCmd.PersistentFlags().StringArrayP("headers", "H", []string{}, "The headers to pass to invocation request.")
	invokeCmd.PersistentFlags().StringArrayP("param", "p", []string{}, "Params to pass to the invocation request.")
}

// validateHttpMethod will return an error if the user has provided an HTTP method that is not supported
func validateHttpMethod(method string) error {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return nil
	}
	return fmt.Errorf("method %s is not supported", method)
}

// invoke a function and get the payload from it.
func invoke(ctx context.Context, opts *invokeOptions) (string, error) {
	// Create an HTTP client that can interact with our Morty backend
	currentCtx := ctx.Value(cliconfig.CurrentCtxKey{}).(*cliconfig.Context)
	cl := httpclient.NewClient(currentCtx.Controller)

	log.Debugf("New invocation request with options: %v", debug.JSON(opts))

	headers := http.Header{}
	// If the caller has passed headers, map them to http.Header
	if opts.Headers != nil {
		for _, header := range opts.Headers {
			splitted := strings.Split(header, ":")
			if len(splitted) != 2 {
				return "", fmt.Errorf("header '%s' is not valid. Please use the correct format: 'Key: Value'", header)
			}
			hKey, hValue := splitted[0], splitted[1]
			headers.Add(hKey, hValue)
		}
	}

	var body io.Reader
	if opts.Body != "" {
		body = bytes.NewBuffer([]byte(opts.Body))
	}

	uri := strings.Replace(invokeFunctionEndpoint, "{name}", opts.FnName, -1)
	uri = strings.Replace(uri, "{version}", opts.FnVersion, -1)

	// If the caller has passed params, add them to url
	if len(opts.Params) > 0 {
		invokeParams := ""
		for _, param := range opts.Params {
			keyValueParam := strings.Split(param, "=")
			if len(keyValueParam) > 1 {
				invokeParams += fmt.Sprintf("%s=%s&", keyValueParam[0], keyValueParam[1])
			} else {
				invokeParams += fmt.Sprintf("%s&", keyValueParam[0])
			}
		}
		uri += fmt.Sprintf("?%s", strings.TrimSuffix(invokeParams, "&"))
	}

	res, err := cl.Generic(ctx, opts.Method, uri, body, headers)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
