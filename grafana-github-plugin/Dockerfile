FROM grafana/grafana

COPY dist /var/lib/grafana/plugins/grafana-github-datasource
COPY ./bin/docker-init.sh /init.sh

# see https://grafana.com/docs/grafana/latest/developers/plugins/sign-a-plugin/#sign-a-plugin
ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS grafana-github-datasource

ENTRYPOINT ["/init.sh"]
