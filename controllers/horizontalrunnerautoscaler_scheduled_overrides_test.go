package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/teambition/rrule-go"
)

func TestCalculateActiveAndUpcomingRecurringPeriods(t *testing.T) {
	type Recurrence struct {
		Start string
		End   string
		Freq  string
		Until string
	}

	type testcase struct {
		Now string

		Recurrence Recurrence

		WantActive   string
		WantUpcoming string
	}

	check := func(t *testing.T, tc testcase) {
		t.Helper()

		_, err := time.Parse(time.RFC3339, "2021-05-08T00:00:00Z")
		if err != nil {
			t.Fatal(err)
		}

		now, err := time.Parse(time.RFC3339, tc.Now)
		if err != nil {
			t.Fatal(err)
		}

		active, upcoming, err := doCalculateActiveAndUpcomingRecurringPeriods(now, tc.Recurrence.Start, tc.Recurrence.End, tc.Recurrence.Freq, tc.Recurrence.Until)
		if err != nil {
			t.Fatal(err)
		}

		if active.String() != tc.WantActive {
			t.Errorf("unexpected active: want %q, got %q", tc.WantActive, active)
		}

		if upcoming.String() != tc.WantUpcoming {
			t.Errorf("unexpected upcoming: want %q, got %q", tc.WantUpcoming, upcoming)
		}
	}

	t.Run("onetime override incoming", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
			},

			Now: "2021-04-30T23:59:59+09:00",

			WantActive:   "",
			WantUpcoming: "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
		})
	})

	t.Run("onetime override started", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
			},

			Now: "2021-05-01T00:00:00+09:00",

			WantActive:   "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("onetime override ending", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
			},

			Now: "2021-05-02T23:59:59+09:00",

			WantActive:   "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("onetime override ended", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
			},

			Now: "2021-05-03T00:00:00+09:00",

			WantActive:   "",
			WantUpcoming: "",
		})
	})

	t.Run("weekly override incoming", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-04-30T23:59:59+09:00",

			WantActive:   "",
			WantUpcoming: "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
		})
	})

	t.Run("weekly override started", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-01T00:00:00+09:00",

			WantActive:   "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
			WantUpcoming: "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
		})
	})

	t.Run("weekly override ending", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-02T23:59:59+09:00",

			WantActive:   "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
			WantUpcoming: "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
		})
	})

	t.Run("weekly override ended", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-03T00:00:00+09:00",

			WantActive:   "",
			WantUpcoming: "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
		})
	})

	t.Run("weekly override reccurrence incoming", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-07T23:59:59+09:00",

			WantActive:   "",
			WantUpcoming: "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
		})
	})

	t.Run("weekly override reccurrence started", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-08T00:00:00+09:00",

			WantActive:   "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
			WantUpcoming: "2021-05-15T00:00:00+09:00-2021-05-17T00:00:00+09:00",
		})
	})

	t.Run("weekly override reccurrence ending", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-09T23:59:59+09:00",

			WantActive:   "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
			WantUpcoming: "2021-05-15T00:00:00+09:00-2021-05-17T00:00:00+09:00",
		})
	})

	t.Run("weekly override reccurrence ended", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-10T00:00:00+09:00",

			WantActive:   "",
			WantUpcoming: "2021-05-15T00:00:00+09:00-2021-05-17T00:00:00+09:00",
		})
	})

	t.Run("weekly override's last reccurrence incoming", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-04-29T23:59:59+09:00",

			WantActive:   "",
			WantUpcoming: "2022-04-30T00:00:00+09:00-2022-05-02T00:00:00+09:00",
		})
	})

	t.Run("weekly override reccurrence started", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-04-30T00:00:00+09:00",

			WantActive:   "2022-04-30T00:00:00+09:00-2022-05-02T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("weekly override reccurrence ending", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-05-01T23:59:59+09:00",

			WantActive:   "2022-04-30T00:00:00+09:00-2022-05-02T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("weekly override reccurrence ended", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-05-02T00:00:00+09:00",

			WantActive:   "",
			WantUpcoming: "",
		})
	})

	t.Run("weekly override repeated forever just starting", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Weekly",
			},

			Now: "2021-05-08T00:00:00+09:00",

			WantActive:   "2021-05-08T00:00:00+09:00-2021-05-10T00:00:00+09:00",
			WantUpcoming: "2021-05-15T00:00:00+09:00-2021-05-17T00:00:00+09:00",
		})
	})

	t.Run("monthly override just starting", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-01T00:00:00+09:00",

			WantActive:   "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
			WantUpcoming: "2021-06-01T00:00:00+09:00-2021-06-03T00:00:00+09:00",
		})
	})

	t.Run("monthly override just recurring", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-06-01T00:00:00+09:00",

			WantActive:   "2021-06-01T00:00:00+09:00-2021-06-03T00:00:00+09:00",
			WantUpcoming: "2021-07-01T00:00:00+09:00-2021-07-03T00:00:00+09:00",
		})
	})

	t.Run("monthly override's last reccurence incoming", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-04-30T23:59:59+09:00",

			WantActive:   "",
			WantUpcoming: "2022-05-01T00:00:00+09:00-2022-05-03T00:00:00+09:00",
		})
	})

	t.Run("monthly override's last reccurence starting", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-05-01T00:00:00+09:00",

			WantActive:   "2022-05-01T00:00:00+09:00-2022-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("monthly override's last reccurence started", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-05-01T00:00:01+09:00",

			WantActive:   "2022-05-01T00:00:00+09:00-2022-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("monthly override's last reccurence ending", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-05-02T23:59:59+09:00",

			WantActive:   "2022-05-01T00:00:00+09:00-2022-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("monthly override's last reccurence ended", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Monthly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2022-05-03T00:00:00+09:00",

			WantActive:   "",
			WantUpcoming: "",
		})
	})

	t.Run("yearly override just starting", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2022-05-01T00:00:00+09:00",
			},

			Now: "2021-05-01T00:00:00+09:00",

			WantActive:   "2021-05-01T00:00:00+09:00-2021-05-03T00:00:00+09:00",
			WantUpcoming: "2022-05-01T00:00:00+09:00-2022-05-03T00:00:00+09:00",
		})
	})

	t.Run("yearly override just recurring", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2023-05-01T00:00:00+09:00",
			},

			Now: "2022-05-01T00:00:00+09:00",

			WantActive:   "2022-05-01T00:00:00+09:00-2022-05-03T00:00:00+09:00",
			WantUpcoming: "2023-05-01T00:00:00+09:00-2023-05-03T00:00:00+09:00",
		})
	})

	t.Run("yearly override's last recurrence incoming", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2023-05-01T00:00:00+09:00",
			},

			Now: "2023-04-30T23:59:59+09:00",

			WantActive:   "",
			WantUpcoming: "2023-05-01T00:00:00+09:00-2023-05-03T00:00:00+09:00",
		})
	})

	t.Run("yearly override's last recurrence starting", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2023-05-01T00:00:00+09:00",
			},

			Now: "2023-05-01T00:00:00+09:00",

			WantActive:   "2023-05-01T00:00:00+09:00-2023-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("yearly override's last recurrence started", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2023-05-01T00:00:00+09:00",
			},

			Now: "2023-05-01T00:00:01+09:00",

			WantActive:   "2023-05-01T00:00:00+09:00-2023-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("yearly override's last recurrence ending", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2023-05-01T00:00:00+09:00",
			},

			Now: "2023-05-02T23:23:59+09:00",

			WantActive:   "2023-05-01T00:00:00+09:00-2023-05-03T00:00:00+09:00",
			WantUpcoming: "",
		})
	})

	t.Run("yearly override's last recurrence ended", func(t *testing.T) {
		t.Helper()

		check(t, testcase{
			Recurrence: Recurrence{
				Start: "2021-05-01T00:00:00+09:00",
				End:   "2021-05-03T00:00:00+09:00",
				Freq:  "Yearly",
				Until: "2023-05-01T00:00:00+09:00",
			},

			Now: "2023-05-03T00:00:00+09:00",

			WantActive:   "",
			WantUpcoming: "",
		})
	})
}

type Period struct {
	StartTime time.Time
	EndTime   time.Time
}

func (r *Period) String() string {
	if r == nil {
		return ""
	}

	return r.StartTime.Format(time.RFC3339) + "-" + r.EndTime.Format(time.RFC3339)
}

func doCalculateActiveAndUpcomingRecurringPeriods(now time.Time, start, end, freq, until string) (*Period, *Period, error) {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, nil, err
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, nil, err
	}

	var untilTime time.Time

	if until != "" {
		ut, err := time.Parse(time.RFC3339, until)
		if err != nil {
			return nil, nil, err
		}

		untilTime = ut
	}

	return CalculateActiveAndUpcomingRecurringPeriods(now, startTime, endTime, freq, untilTime)
}

func CalculateActiveAndUpcomingRecurringPeriods(now, startTime, endTime time.Time, freq string, untilTime time.Time) (*Period, *Period, error) {
	var freqValue rrule.Frequency

	var freqDurationDay int
	var freqDurationMonth int
	var freqDurationYear int

	switch freq {
	case "Daily":
		freqValue = rrule.DAILY
		freqDurationDay = 1
	case "Weekly":
		freqValue = rrule.WEEKLY
		freqDurationDay = 7
	case "Monthly":
		freqValue = rrule.MONTHLY
		freqDurationMonth = 1
	case "Yearly":
		freqValue = rrule.YEARLY
		freqDurationYear = 1
	case "":
		if now.Before(startTime) {
			return nil, &Period{StartTime: startTime, EndTime: endTime}, nil
		}

		if now.Before(endTime) {
			return &Period{StartTime: startTime, EndTime: endTime}, nil, nil
		}

		return nil, nil, nil
	default:
		return nil, nil, fmt.Errorf(`invalid freq %q: It must be one of "Daily", "Weekly", "Monthly", and "Yearly"`, freq)
	}

	freqDurationLater := time.Date(
		now.Year()+freqDurationYear,
		time.Month(int(now.Month())+freqDurationMonth),
		now.Day()+freqDurationDay,
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location(),
	)

	freqDuration := freqDurationLater.Sub(now)

	overrideDuration := endTime.Sub(startTime)
	if overrideDuration > freqDuration {
		return nil, nil, fmt.Errorf("override's duration %s must be equal to sor shorter than the duration implied by freq %q (%s)", overrideDuration, freq, freqDuration)
	}

	rrule, err := rrule.NewRRule(rrule.ROption{
		Freq:    freqValue,
		Dtstart: startTime,
		Until:   untilTime,
	})
	if err != nil {
		return nil, nil, err
	}

	overrideDurationBefore := now.Add(-overrideDuration + 1)
	activeOverrideStarts := rrule.Between(overrideDurationBefore, now, true)

	var active *Period

	if len(activeOverrideStarts) > 1 {
		return nil, nil, fmt.Errorf("[bug] unexpted number of active overrides found: %v", activeOverrideStarts)
	} else if len(activeOverrideStarts) == 1 {
		active = &Period{
			StartTime: activeOverrideStarts[0],
			EndTime:   activeOverrideStarts[0].Add(overrideDuration),
		}
	}

	oneSecondLater := now.Add(1)
	upcomingOverrideStarts := rrule.Between(oneSecondLater, freqDurationLater, true)

	var next *Period

	if len(upcomingOverrideStarts) > 0 {
		next = &Period{
			StartTime: upcomingOverrideStarts[0],
			EndTime:   upcomingOverrideStarts[0].Add(overrideDuration),
		}
	}

	return active, next, nil
}
