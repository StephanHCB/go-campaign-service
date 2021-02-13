package acceptance

import (
	"github.com/StephanHCB/go-campaign-service/web"
	"github.com/StephanHCB/go-campaign-service/web/controller/healthctl"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
	"time"
)

// ------------------------------------------
// acceptance tests for the timeout resource
// ------------------------------------------

func TestTimeoutResponse(t *testing.T) {
	// make this test run in reasonable time
	healthctl.SleepTime = 5 * time.Millisecond
	web.RequestTimeout = 25 * time.Millisecond

	tstSetup(tstValidConfigurationPath)
	defer tstShutdown()

	Convey("Given an authenticated user", t, func() {
		token := tstAuthUser()

		Convey("When they make a request that runs too long", func() {
			response, err := tstPerformGet("/timeout", token)

			Convey("Then a timeout occurs", func() {
				So(err, ShouldEqual, nil)
				So(response.status, ShouldEqual, http.StatusGatewayTimeout)
			})
		})
	})
}
