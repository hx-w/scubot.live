<template>
    <div>
        <a-empty :description="cuid">在做了在做了...</a-empty>
    </div>
</template>
<script>
export default {
    props: ['uid'],
    data() {
        return {
            cuid : this.uid
        }
    },
    watch: {
        uid(newVal) {
            this.cuid = newVal
        }
    },
    methods: {
    }
}
</script>
