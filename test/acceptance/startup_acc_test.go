package acceptance

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestStartupWithValidConfig_ShouldBeHealthy(t *testing.T) {
	Convey("Given a valid configuration", t, func() {
		configAndSecretsPath := tstValidConfigurationPath

		Convey("When the application is started", func() {
			tstSetup(configAndSecretsPath)
			defer tstShutdown()

			Convey("Then it reports as healthy", func() {
				response, err := tstPerformGet("/health", tstUnauthenticated())

				So(err, ShouldEqual, nil)
				So(response.status, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestStartupWithInvalidConfig_ShouldNotRespond(t *testing.T) {
	Convey("Given an invalid configuration", t, func() {
		configAndSecretsPath := tstInvalidConfigurationPath

		Convey("When the application is started", func() {
			tstSetup(configAndSecretsPath)
			defer tstShutdown()

			Convey("Then server startup is aborted with the expected error", func() {
				So(failures[0].Error(), ShouldEqual, "Fatal error: configuration value for key server.port is not in range 1024..65535")
			})

			Convey("And it does not respond to health inquiries", func() {
				_, err := tstPerformGet("/health", tstUnauthenticated())

				So(err, ShouldNotBeNil)
			})
		})
	})
}
