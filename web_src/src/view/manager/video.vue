<template>
  <div
    class="device_container"
    id="device_container"
  >
    <!-- <video tabindex="-1"></video> -->
  </div>
</template>
<script>
import { startStream, stopStream } from '@/api/device'
import ChimeePlayer from 'chimee-player';

export default {
  data () {
    return {
      params: {
        channel: '',
        deviceID: ''
      }

    }
  },
  computed: {
  },
  created () {
    if (this.$route.query.deviceID) {
      this.params = {
        deviceID: this.$route.query.deviceID,
        channelID: this.$route.query.channel
      }
    }
  },
  mounted () {
    this.startLive()
  },
  methods: {
    startLive () {
      startStream(this.params)
        .then(res => {
          console.log(res)
          if (res.errCode == 200) {
            new ChimeePlayer({
              wrapper: '#device_container',
              src: res.sessionURL.flv,
              box: 'flv',
              isLive: true,
              autoplay: true,
              controls: true
            });
          } else {
            this.$message({
              message: res.errMsg,
              type: 'error'
            });
          }

        })
    },
    stopLive () {
      stopStream(this.params)
        .then(res => {
          console.log(res)
        })
    }
  }
}
</script>
<style lang="css">
@import url('./chimee-player.browser.css');
.device_container {
  height: 100%;
  width: 100%;
  position: relative;
  display: block;
}
.chimee-container {
  position: relative;
}
</style>