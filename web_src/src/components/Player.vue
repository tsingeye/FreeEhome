<template>
  <!-- 根据 monitor 父组件传递过来的数据，判断容器是否需要高亮显示 -->
  <div :class="['player-container', { 'active-index': activeIndex === index }]">
    <div class="player-title" v-if="id !== ''">
      <span> {{ title }}</span>
      <i class="el-icon-close" @click="handleClose" />
    </div>
    <EasyPlayer
      :videoUrl="url"
      live
      :fluent="fluent"
      :autoplay="autoplay"
      :aspect="aspect"
      stretch
    ></EasyPlayer>
  </div>
</template>

<script>
import EasyPlayer from '@easydarwin/easyplayer'
export default {
  name: 'Player',
  components: { EasyPlayer },
  props: {
    activeIndex: {
      type: Number,
      default: 0
    },
    index: {
      type: Number,
      default: 0
    },
    id: {
      type: String,
      default: ''
    },
    title: {
      type: String,
      default: ''
    },
    url: {
      type: String,
      default: ''
    },
    aspect: {
      type: String,
      default: '100%'
    }
  },
  data() {
    return {
      fluent: true,
      autoplay: true
    }
  },
  methods: {
    handleClose() {
      this.$emit('close', this.index)
    }
  }
}
</script>

<style></style>
