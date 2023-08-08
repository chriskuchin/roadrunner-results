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
        ports = ["http"]

        entrypoint = []
        volumes = [
          "/etc/localtime:/etc/localtime:ro",
          "/var/lib/rslts:/rslts"
        ]
      }

      vault {
        policies    = ["cloudflare-r2"]
        change_mode = "restart"
      }

      template {
        data        = <<-EOH
        {{with secret "kv/rslts/api" }}
        API_TOKEN={{.Data.data.token}}
        {{end}}
        EOH
        destination = "local/api.env"
        env         = true
      }

      service {
        name = "rslts"
        port = "http"
        tags = [
          "traefik.enable=true",
          "traefik.http.routers.rslts_http.entrypoints=http",
          "traefik.http.routers.rslts_http.rule=Host(`rslts.home.cksuperman.com`)",
          "traefik.http.routers.rslts_http.middlewares=rslts-redirect@consulcatalog",
          "traefik.http.middlewares.rslts-redirect.redirectscheme.scheme=https",
          "traefik.http.middlewares.rslts-redirect.redirectscheme.permanent=true",
          "traefik.http.routers.rslts.entrypoints=https",
          "traefik.http.routers.rslts.rule=Host(`rslts.home.cksuperman.com`)",
          "traefik.http.routers.rslts.tls.certresolver=cloudflare",
          "traefik.http.routers.rslts.tls.domains[0].main=rslts.home.cksuperman.com",
          "wayfinder.domain=rslts.home.cksuperman.com",
        ]
      }

      resources {
        cpu    = 100
        memory = 256
      }
    }
  }

}