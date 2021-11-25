<template>
  <div>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column type="index" label="#"> </el-table-column>
      <el-table-column prop="deviceID" label="设备ID"> </el-table-column>
      <el-table-column prop="channelID" label="通道ID"> </el-table-column>
      <el-table-column prop="channelName" label="通道名称"> </el-table-column>
      <el-table-column prop="status" label="设备状态"> </el-table-column>
      <el-table-column prop="updatedAt" label="更新时间"> </el-table-column>
      <el-table-column prop="createdAt" label="创建时间"> </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button type="primary" @click="handleChannel(scope.row)"
            >播放</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <el-dialog
      title="播放"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
    >
      <player :url="playUrl" aspect="16:9" />
    </el-dialog>
  </div>
</template>

<script>
import { list } from '@/api/channels'
import { stream } from '@/api/stream'
import Player from '@/components/Player.vue'
export default {
  components: {
    Player
  },
  data() {
    return {
      tableData: [],
      dialogVisible: false,
      playUrl: ''
    }
  },
  created() {
    this.getChannelList()
  },
  methods: {
    async getChannelList() {
      const token = window.localStorage.getItem('token')
      const deviceId = this.$route.query.deviceId
      const { data } = await list(token, deviceId)
      if (data.errCode !== 200) return this.$notify.error('请求失败！')
      this.tableData = data.channelList
    },
    async handleChannel(row) {
      const token = window.localStorage.getItem('token')
      const { data } = await stream(token, row.channelID)
      if (data.errCode !== 200) return this.$notify.error('请求失败！')

      console.log(data)
      this.playUrl = data.sessionURL.hls
      // this.playUrl = 'http://cctvalih5ca.v.myalicdn.com/live/cctv1_2/index.m3u8'
      this.dialogVisible = true
    }
  }
}
</script>
