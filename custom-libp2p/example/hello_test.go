package example

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomLibP2P_GetHelloMessage(t *testing.T) {
	testerUserName := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis blandit massa enim, non semper leo commodo vitae. Nullam condimentum velit et diam vestibulum, nec mattis dui suscipit. Nullam eu viverra velit. Phasellus pretium lectus neque, sed vestibulum lorem viverra sit amet. Proin vulputate magna et mauris hendrerit, sed porttitor ante pulvinar. Praesent augue ipsum, lobortis vel justo bibendum, rhoncus sollicitudin diam. Suspendisse vehicula dolor sed lorem ornare lacinia. "
	output, err := GetHelloMessage(testerUserName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello "+testerUserName+" this is a dummy function!", output)
	t.Log("Output: " + output)
}
