<template>
  <div class="device_container">
    <el-form
      :label-position="'right'"
      :model="ruleForm"
      ref="ruleForm"
      :inline="true"
    >
      <el-form-item
        label="状态:"
        prop='status'
      >
        <el-select
          v-model="ruleForm.status"
          placeholder="请选择"
          @change='selectChange'
        >
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          :loading="device.loading"
          @click="search"
        >搜索</el-button>

      </el-form-item>
    </el-form>
    <div>
      <el-table
        :data="device.deviceConfig.deviceList"
        border
        style="width: 100%"
      >
        <el-table-column
          prop="deviceID"
          label="设备ID"
          min-width="150"
        >
          <template slot-scope="scope">
            <span class="device_id" @click='toChannel(scope.row.deviceID)'>{{scope.row.deviceID}}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="deviceIP"
          label="设备IP"
          min-width="150"
        >
        </el-table-column>
        <el-table-column
          prop="deviceName"
          label="设备名"
          min-width="150"
        >
        </el-table-column>
        <el-table-column
          prop="serialNumber"
          label="设备序列号"
          min-width="150"
        >
        </el-table-column>
        <el-table-column
          prop="status"
          label="状态"
          min-width="150"
        >
          <template slot-scope="scope">
            <span>{{scope.row.status == 'ON' ? '在线' : '离线'}}</span>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination_view">
        <el-pagination
          background
          layout="prev, pager, next"
          :page-size='10'
          :total="device.deviceConfig.totalCount"
          :current-page='ruleForm.page'
          @current-change='pageChange'
        >
        </el-pagination>
      </div>
    </div>
  </div>
</template>
<script>
import { deviceList } from '@/api/device'
import { createNamespacedHelpers,mapState } from "vuex";
let { mapActions } = createNamespacedHelpers("device");
export default {
  data () {
    return {
      options: [
        {
          value: 'ON',
          label: '在线'
        },
        {
          value: 'OFF',
          label: '离线'
        }
      ],
      ruleForm: {
        status: '',
        page: 1,
        limit: 10
      },
      loading: false,
    }
  },
  computed: {
      ...mapState(['device'])
  },
  created () {
    this.getDevice()
  },
  methods: {
    ...mapActions(['getDeviceList']),
    selectChange (event) {
        // this.ruleForm.page=1
     },
     pageChange(event){
         this.ruleForm.page=event
         this.getDevice()
     },
     search(){
        this.ruleForm.page=1
        this.getDevice()
     },
    getDevice () {
        this.getDeviceList(this.ruleForm)
    },
    toChannel(deviceID){
        console.log(deviceID)
        this.$router.push('/manager/channel?deviceID='+deviceID)
    }
  }
}
</script>
<style lang="css">
.device_container .el-form--inline .el-form-item {
  float: left;
}
.device_id {
  color: dodgerblue;
  cursor: pointer;
}
.pagination_view {
  margin-top: 10px;
}
</style>