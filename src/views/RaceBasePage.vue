<template>
    <div>
        <breadcrumb :paths="paths"></breadcrumb>
        <router-view></router-view>
    </div>
</template>

<script>
import Breadcrumb from "../components/Breadcrumb.vue";

export default {
    components: {
        breadcrumb: Breadcrumb,
    },
    data: function () {
        return {
            race: {
                name: "",
                owner: "",
                id: "",
            }
        }
    },
    mounted: function () {
        this.getRaceInfo()
    },
    methods: {
        getRaceInfo: function () {
            fetch("/api/v1/races/" + this.$route.params.raceId, {
                method: "GET"
            }).then(res => {
                return res.json()
            }).then(res => {
                console.log(res)
                this.race.name = res.name
            })
        },
    },
    computed: {
        paths: function () {
            var breadcrumb = [{
                path: "/races",
                display: "Races"
            }]
            var path = "/races"
            var pathSegments = this.$route.path.replace("/races/", "").split("/")
            for (let i = 0; i < pathSegments.length; i++) {
                var display = pathSegments[i]
                if (i == 0 && this.race.name != "") {
                    display = this.race.name
                }
                path += "/" + pathSegments[i]
                breadcrumb.push({
                    path: path,
                    display: display
                })
            }

            return breadcrumb;
        },
    },
};
</script>