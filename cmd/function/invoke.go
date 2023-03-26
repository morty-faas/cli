package function

import (
	"encoding/json"
	"fmt"
	"morty/cliconfig"
	"morty/client/gateway"
	"net/http"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke NAME",
	Short: "Invoke a function",
	Long:  `Invoke a function using default options or choose HTTP method, body, headers etc.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context().Value(cliconfig.CurrentCtxKey{}).(*cliconfig.Context)

		// Safe call, validation is performed by cobra.ExactArgs(1) above
		name := args[0]

		// Extract the method from the command
		method, _ := cmd.Flags().GetString("method")
		method = strings.ToUpper(method)
		if err := validateHttpMethod(method); err != nil {
			return err
		}

		data, _ := cmd.Flags().GetString("data")
		headers, _ := cmd.Flags().GetStringArray("headers")
		jsonPathQuery, _ := cmd.Flags().GetString("query")

		gw := gateway.NewClient(ctx.Gateway)

		opts := &gateway.InvokeFnRequest{
			FnName:  name,
			Method:  method,
			Body:    data,
			Headers: headers,
		}

		response, err := gw.InvokeFn(cmd.Context(), opts)
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
			fmt.Println(string(by))
			return nil
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	invokeCmd.PersistentFlags().StringP("method", "X", "GET", "The HTTP method to use to invoke the request. Valid values are: GET, POST, PUT, PATCH, DELETE")
	invokeCmd.PersistentFlags().StringP("data", "d", "", "The body to pass to the invocation request.")
	invokeCmd.PersistentFlags().StringP("query", "q", "$", "A valid JSON Path expression to execute on the function response. If the function output isn't a JSON, the flag will have no effect.")
	invokeCmd.PersistentFlags().StringArrayP("headers", "H", []string{}, "The headers to pass to invocation request.")
}

// validateHttpMethod will return an error if the user has provided an HTTP method that is not supported
func validateHttpMethod(method string) error {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return nil
	}
	return fmt.Errorf("method %s is not supported", method)
}
