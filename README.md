# Virgin Hub Prometheus Exporter

![Docker Build](https://github.com/mowaxuk/virgin-hub-exporter/actions/workflows/docker-build.yml/badge.svg)

A lightweight Prometheus exporter for reverse-engineered Virgin Media Hub telemetry endpoints.  
Built with Go, containerized with Docker, and published for public use.

## 🚀 Features

- Exposes undocumented metrics from Virgin Media Hub routers
- Outputs metrics in Prometheus format
- Compatible with Prometheus and Grafana
- Includes a prebuilt Grafana dashboard
- Dockerized for easy deployment
- CI/CD pipeline via GitHub Actions

## 📦 Docker Usage

Run the exporter as a container:

```bash
docker run -d -p 9111:9111 mowaxuk/virgin-hub-exporter:latest
```

Metrics will be available at: `http://localhost:9111/metrics`

## 📊 Grafana Dashboard

Import the included dashboard JSON to visualize metrics:

`grafana/virgin-hub-dashboard.json`

This dashboard includes panels for:

- Downstream and upstream signal strength
- Connection status
- Uptime
- Channel metrics
- Modulation and error rates

## 🔧 Prometheus Scrape Config

Add this to your Prometheus configuration:

```yaml
- job_name: 'virgin-hub'
  static_configs:
    - targets: ['localhost:9111']
```

Replace `localhost` with your host IP if running remotely.

## 🧠 Reverse Engineering Notes

This exporter was reverse engineered using a Virgin Media Hub, and pinpointing telemetry endpoints by **Mowaxuk (Wayne)** (GitHub: `mowaxuk`)
It exposes undocumented downstream and upstream signal metrics from consumer-grade routers manufactured by **Technicolor** and **Hitron**.  
The metrics were extracted by inspecting the router’s web interface and HTTP endpoints, and translated into Prometheus format for sysadmin-grade observability.  
Built with Go, containerized with Docker, and published for public use.

## 📜 License

MIT License. See `LICENSE` file for details.
