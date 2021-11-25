<template>
  <div>
    <el-tree
      :data="data"
      :props="defaultProps"
      @node-click="handleNodeClick"
    ></el-tree>
  </div>
</template>
<script>
export default {
  data() {
    return {
      data: [
        {
          label: 'test1',
          children: [
            {
              label: 'test1',
              name: 'test1',
              url: 'http://cctvalih5ca.v.myalicdn.com/live/cctv1_2/index.m3u8'
            }
          ]
        },
        {
          label: 'test2',
          children: [
            {
              label: 'test2',
              name: 'test2',
              url: 'http://220.161.87.62:8800/hls/1/index.m3u8'
            }
          ]
        }
      ],
      defaultProps: {
        children: 'children',
        label: 'label'
      }
    }
  },
  methods: {
    handleNodeClick(data) {
      this.$emit('clickInfo', data)
    }
  }
}
</script>
