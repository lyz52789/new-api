package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "new_api"

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status", "tag"},
	)

	HTTPRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30},
		},
		[]string{"method", "route", "status", "tag"},
	)

	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "http_requests_in_flight",
			Help:      "Current number of in-flight HTTP requests.",
		},
		[]string{"tag"},
	)

	TopupSuccessTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "topup_success_total",
			Help:      "Total number of successful topups.",
		},
		[]string{"payment_method", "source"},
	)

	TopupSuccessMoneyCentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "topup_success_money_cents_total",
			Help:      "Total successful topup payment money in cents.",
		},
		[]string{"payment_method", "source"},
	)

	AffiliateCommissionCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "affiliate_commission_created_total",
			Help:      "Total number of affiliate commission records created.",
		},
		[]string{"payment_method", "status"},
	)

	AffiliateCommissionAmountCentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "affiliate_commission_amount_cents_total",
			Help:      "Total affiliate commission amount in cents.",
		},
		[]string{"payment_method", "status"},
	)
)

func ObserveHTTPRequest(method, route string, status int, tag string, duration time.Duration) {
	if route == "" {
		route = "unmatched"
	}
	if tag == "" {
		tag = "web"
	}
	statusLabel := strconv.Itoa(status)
	HTTPRequestsTotal.WithLabelValues(method, route, statusLabel, tag).Inc()
	HTTPRequestDurationSeconds.WithLabelValues(method, route, statusLabel, tag).Observe(duration.Seconds())
}

func IncHTTPInFlight(tag string) {
	if tag == "" {
		tag = "web"
	}
	HTTPRequestsInFlight.WithLabelValues(tag).Inc()
}

func DecHTTPInFlight(tag string) {
	if tag == "" {
		tag = "web"
	}
	HTTPRequestsInFlight.WithLabelValues(tag).Dec()
}

func RecordTopupSuccess(paymentMethod string, source string, moneyCents int64) {
	if paymentMethod == "" {
		paymentMethod = "unknown"
	}
	if source == "" {
		source = "unknown"
	}
	TopupSuccessTotal.WithLabelValues(paymentMethod, source).Inc()
	if moneyCents > 0 {
		TopupSuccessMoneyCentsTotal.WithLabelValues(paymentMethod, source).Add(float64(moneyCents))
	}
}

func RecordAffiliateCommissionCreated(paymentMethod string, status string, amountCents int64) {
	if paymentMethod == "" {
		paymentMethod = "unknown"
	}
	if status == "" {
		status = "unknown"
	}
	AffiliateCommissionCreatedTotal.WithLabelValues(paymentMethod, status).Inc()
	if amountCents > 0 {
		AffiliateCommissionAmountCentsTotal.WithLabelValues(paymentMethod, status).Add(float64(amountCents))
	}
}
