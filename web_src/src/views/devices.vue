<template>
  <el-table :data="tableData" style="width: 100%">
    <el-table-column type="index" label="#"> </el-table-column>
    <el-table-column prop="deviceID" label="设备ID"> </el-table-column>
    <el-table-column prop="deviceIP" label="设备IP"> </el-table-column>
    <el-table-column prop="deviceName" label="设备名称"> </el-table-column>
    <el-table-column prop="serialNumber" label="序号/编号"> </el-table-column>
    <el-table-column prop="status" label="设备状态"> </el-table-column>
    <el-table-column prop="updatedAt" label="更新时间"> </el-table-column>
    <el-table-column prop="createdAt" label="创建时间"> </el-table-column>
    <el-table-column label="操作">
      <template slot-scope="scope">
        <el-button type="primary" @click="handleChannel(scope.row)"
          >查看通道</el-button
        >
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
import { list } from '@/api/devices'
export default {
  name: 'Devices',
  data() {
    return {
      tableData: []
    }
  },
  created() {
    this.getDeviceList()
  },
  methods: {
    async getDeviceList() {
      const token = window.localStorage.getItem('token')
      const { data } = await list(token)
      if (data.errCode !== 200) return this.$notify.error('请求失败！')
      this.tableData = data.deviceList
    },
    handleChannel(row) {
      this.$router.push({
        path: '/channels',
        query: {
          deviceId: row.deviceID
        }
      })
    }
  }
}
</script>
