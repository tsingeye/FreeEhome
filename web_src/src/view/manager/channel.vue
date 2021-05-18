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
        :data="device.channelConfig.channelList"
        border
        style="width: 100%"
      >
        <el-table-column
          prop="channelID"
          label="通道ID"
          min-width="150"
        >
          <template slot-scope="scope">
            <span
              class="device_id"
              @click='toChannel(scope.row.channelID)'
            >{{scope.row.channelID}}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="deviceID"
          label="设备ID"
          min-width="150"
        >
        </el-table-column>

        <el-table-column
          prop="channelName"
          label="通道名称"
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
          :total="device.channelConfig.totalCount"
          :current-page='ruleForm.page'
          @current-change='pageChange'
        >
        </el-pagination>
      </div>
    </div>
    <el-dialog
      title="查看"
      :visible.sync="dialogVisible"
      width="70%"
      :before-close='closeDialog'
    >
      <div
        class="device_container"
        id="device_container"
      >
      </div>
      <span
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="closeDialog">取 消</el-button>
        <el-button
          type="primary"
          @click="closeDialog"
        >确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import { deviceList, startStream, stopStream } from '@/api/device'
import ChimeePlayer from 'chimee-player';
// import LivePlayer from '@liveqing/liveplayer'
import { createNamespacedHelpers, mapState } from "vuex";
let loadingFull
let { mapActions } = createNamespacedHelpers("device");

var playerFlv = null
export default {
  data () {
    return {
      dialogVisible: false,
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
        limit: 10,
        deviceID: ''
      },
      loading: false,
      liveForm: {},
    }
  },
  computed: {
    ...mapState(['device'])
  },
  created () {
    if (this.$route.query.deviceID) {
      this.ruleForm.deviceID = this.$route.query.deviceID
    }
    this.getDevice()
  },
  methods: {
    ...mapActions(['getChannelList']),
    selectChange (event) {
      // this.ruleForm.page=1
    },
    pageChange (event) {
      this.ruleForm.page = event
      this.getDevice()
    },
    search () {
      this.ruleForm.page = 1
      this.getDevice()
    },
    getDevice () {
      this.getChannelList(this.ruleForm)
    },
    toChannel (channel) {
      // this.$router.push('/manager/video?channel=' + channel + '&deviceID=' + this.ruleForm.deviceID)
      this.startLive({
        channelID: channel,
        deviceID: this.ruleForm.deviceID
      })
    },
    closeDialog () {
      this.dialogVisible = false
      this.stopLive()
    },
    stopLive () {
      console.log(this.liveForm, 'this.liveForm.')
      if (this.liveForm.channelID) {
        stopStream(this.liveForm)
          .then(res => {
            console.log(res)
          })
      }
      if (playerFlv) {
        playerFlv.stopLoad()
      }

    },
    startLive (params) {
      loadingFull = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      startStream(params)
        .then(res => {
          loadingFull.close()
          console.log(res)
          if (res.errCode == 200) {
            this.liveForm = params
            this.dialogVisible = true
            setTimeout(() => {
              playerFlv = new ChimeePlayer({
                wrapper: '#device_container',
                src: res.sessionURL.flv,
                box: 'flv',
                isLive: true,
                autoplay: true,
                controls: true,
                muted:true
              });
            }, 200);

          } else {
            this.liveForm = {}
            this.$message({
              message: res.errMsg,
              type: 'error'
            });
          }

        })
    },
  }
}
</script>
<style lang="css">
@import url('./chimee-player.browser.css');
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