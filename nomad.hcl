job "rslts" {
  datacenters = ["home1"]

  type = "service"

  group "rslts" {
    count = 1

    network {
      port "http" {
        to = 3030
      }
    }

    task "rslts" {
      driver = "docker"

      config {
        image = "ghcr.io/chriskuchin/roadrunner-results:main"
        force_pull = true
        ports = ["http"]

        entrypoint = []
        volumes = [
          "/etc/localtime:/etc/localtime:ro",
          "/var/lib/rslts:/rslts"
        ]
      }

      vault {
        policies    = ["rslts"]
        change_mode = "restart"
      }

      template {
        data        = <<-EOH
{{with secret "kv/rslts/api" }}
API_TOKEN={{.Data.data.token}}
{{end}}
        EOH
        destination = "secrets/api.env"
        env         = true
      }

      template {
        data = <<-EOH
{{with secret "kv/cloudflare/r2/rslts"}}
R2_ACCOUNT_ID={{.Data.data.account}}
R2_ACCESS_KEY_ID={{.Data.data.access_key}}
R2_SECRET_ACCESS_KEY={{.Data.data.secret_key}}
{{end}}
        EOH

        destination = "secrets/cloudflare.env"
        env = true
      }

      template {
        data        = <<-EOH
{{with secret "kv/rslts/auth" }}
{{.Data.data | toJSON }}
{{end}}
        EOH
        destination = "secrets/auth.json"
        env         = false
      }

      service {
        name = "rslts"
        port = "http"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.rslts.entrypoints=public_https",
          "traefik.http.routers.rslts.rule=Host(`rslts.run`)",
          "traefik.http.routers.rslts.tls.certresolver=http",
          "traefik.http.routers.rslts.tls.domains[0].main=rslts.run",
          "wayfinder.domain=rslts.cksuperman.com",
          "wayfinder.public=true",
        ]
      }

      resources {
        cpu    = 100
        memory = 256
      }
    }
  }

}