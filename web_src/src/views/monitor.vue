<template>
  <div class="container">
    <div class="layout-main-sidebar">
      <TreeList @clickInfo="onClickInfo" />
    </div>
    <div class="layout-main-container">
      <div class="monitor-screen-top">
        <div
          v-for="(item1, index1) in layoutPlayerList"
          :key="index1"
          :class="[
            'layout-screen-btn',
            item1.containerClass,
            { active: item1.id === screenTypeActive }
          ]"
          @click="handleScreenType(item1)"
        >
          <div
            v-for="(item2, index2) in item1.children"
            :key="index2"
            :class="item2.itemClass"
          ></div>
        </div>
        <i class="el-icon-full-screen" @click="toggle()"></i>
      </div>
      <fullscreen
        :fullscreen.sync="fullscreen"
        class="monitor-screen-container"
      >
        <div :class="['layout-screen-container', playList.containerClass]">
          <div
            v-for="(item, index) in playList.children"
            :key="index"
            :class="item.itemClass"
            @click="onActiveIndex(index + 1)"
          >
            <Player
              :active-index="activeIndex"
              :index="index + 1"
              :id="item.id"
              :title="item.title"
              :url="item.url"
              @close="onClose"
            />
          </div>
        </div>
      </fullscreen>
    </div>
  </div>
</template>

<script>
import TreeList from '@/components/TreeList.vue'
import Player from '@/components/Player.vue'
export default {
  name: 'Monitor',
  components: {
    Player,
    TreeList
  },
  data() {
    return {
      screenTypeActive: 3, // 当前显示模式
      activeNum: 6, // 显示数量
      fullscreen: false, // 是否全屏
      // 当前播放对象
      playList: {
        id: 3,
        type: 6,
        containerClass: 'layout-screen_3',
        children: [
          { id: '', title: '', url: '', itemClass: 'screen-item_6-1' },
          { id: '', title: '', url: '' },
          { id: '', title: '', url: '' },
          { id: '', title: '', url: '' },
          { id: '', title: '', url: '' },
          { id: '', title: '', url: '' }
        ]
      },
      // 预定义所有对象
      layoutPlayerList: [
        {
          id: 1,
          type: 1,
          containerClass: 'layout-screen_1',
          children: [{ id: '', title: '', url: '' }]
        },
        {
          id: 2,
          type: 4,
          containerClass: 'layout-screen_2',
          children: [
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' }
          ]
        },
        {
          id: 3,
          type: 6,
          containerClass: 'layout-screen_3',
          children: [
            { id: '', title: '', url: '', itemClass: 'screen-item_6-1' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' }
          ]
        },
        // {
        //   id: 4,
        //   type: 8,
        //   containerClass: 'layout-screen_4',
        //   children: [
        //     { id: '', title: '', url: '', itemClass: 'screen-item_8-1' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' }
        //   ]
        // },
        {
          id: 5,
          type: 9,
          containerClass: 'layout-screen_3',
          children: [
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' }
          ]
        },
        // {
        //   id: 6,
        //   type: 10,
        //   containerClass: 'layout-screen_5',
        //   children: [
        //     { id: '', title: '', url: '', itemClass: 'screen-item_101-1' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' }
        //   ]
        // },
        // {
        //   id: 7,
        //   type: 10,
        //   containerClass: 'layout-screen_4',
        //   children: [
        //     { id: '', title: '', url: '', itemClass: 'screen-item_102-1' },
        //     { id: '', title: '', url: '', itemClass: 'screen-item_102-2' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' }
        //   ]
        // },
        // {
        //   id: 8,
        //   type: 13,
        //   containerClass: 'layout-screen_5',
        //   children: [
        //     { id: '', title: '', url: '', itemClass: 'screen-item_13-1' },
        //     { id: '', title: '', url: '', itemClass: 'screen-item_13-2' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '', itemClass: 'screen-item_13-5' },
        //     { id: '', title: '', url: '', itemClass: 'screen-item_13-6' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' },
        //     { id: '', title: '', url: '' }
        //   ]
        // },
        {
          id: 9,
          type: 16,
          containerClass: 'layout-screen_4',
          children: [
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' },
            { id: '', title: '', url: '' }
          ]
        }
      ],
      activeIndex: 0 // 当前激活播放器
    }
  },
  methods: {
    // 点击头部按钮，如果点击的按钮，是当前显示的模式 return
    // 否则获取当前点击的 id，根据 id 匹配到预定义的数据，赋值给 playList
    // 设置当前显示模式
    handleScreenType(item) {
      if (this.screenTypeActive === item.id) return
      this.screenTypeActive = item.id
      const data = this.layoutPlayerList.find(
        item => item.id === this.screenTypeActive
      )
      this.playList = data
      this.activeNum = item.type
    },
    // 是否全屏显示
    toggle() {
      this.fullscreen = !this.fullscreen
    },
    // TreeList 组件传递过来的方法，点击获取节点数据
    // activeIndex 为 0 表示，没有预选播放位置
    // 设置当前播放对象的 id，title，url 等需要配置的属性，传递给 Player 组件
    // 因为 activeIndex 从 0 开始 ，如果 activeIndex >= activeNum 表示已经到了最后一个
    // 否则 activeIndex ++ ，播放位置到下一个
    onClickInfo(val) {
      if (this.activeIndex === 0) this.activeIndex = 1
      this.playList.children[this.activeIndex - 1].id = this.activeIndex + ''
      this.playList.children[this.activeIndex - 1].title = val.name
      this.playList.children[this.activeIndex - 1].url = val.url
      if (this.activeIndex >= this.activeNum) {
        this.activeIndex = 1
      } else {
        this.activeIndex++
      }
    },
    // 设置 activeIndex 为当前播放的 index
    // 目的：高亮显示容器、控制播放位置
    onActiveIndex(index) {
      this.activeIndex = index
    },
    // Player 组件传递过来的方法，点击关闭视频
    onClose(val) {
      this.playList.children[val - 1].id = ''
      this.playList.children[val - 1].title = ''
      this.playList.children[val - 1].url = ''
    }
  }
}
</script>

<style lang="less" scoped>
.container {
  display: flex;
  height: 100%;
}
</style>
