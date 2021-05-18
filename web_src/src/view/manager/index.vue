<template>
  <div>
    <el-container>
      <el-aside width="200px">
        <div class="logo">logo</div>
        <el-col :span="24">
          <el-menu
            class="el-menu-vertical-demo"
            @open="handleOpen"
            @close="handleClose"
            router
            :default-active="$route.path"
            active-text-color="#000"
          >
            <el-menu-item index="/manager/device">
              <i class="el-icon-menu"></i>
              <span slot="title">设备列表</span>
            </el-menu-item>
          </el-menu>
        </el-col>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header_menu">
            <span
              class="logout"
              @click="logout"
            >退出</span>
          </div>
        </el-header>
        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>
<script>
import axios from '@/utils/request'
export default {
  methods: {
    handleOpen (key, keyPath) {
      console.log(key, keyPath);
    },
    handleClose (key, keyPath) {
      console.log(key, keyPath);
    },
    logout () {
      const that = this
      axios.get('/api/v1/system/logout', {})
        .then(res => {
          localStorage.removeItem('authCode')
          that.$router.push('/login')

        })
    }
  }
}
</script>
<style lang="css">
.el-header {
  background-color: #b3c0d1;
  color: #333;
  text-align: center;
  line-height: 60px;
}
.el-aside {
  background-color: #d3dce6;
  color: #333;
  text-align: center;
  height: 100%;
  float: left;
}
.el-main {
  background-color: #e9eef3;
  color: #333;
  text-align: center;
  height: calc(100vh - 60px);
}
.el-container {
  height: 100vh;
}
.logo {
  height: 60px;
  width: 200px;
}
.el-menu-vertical-demo {
  height: calc(100vh - 60px);
}
.el-menu-item.is-active {
  color: #6681fa;
  background-color: #eaeeff !important;
}
.logout {
  float: right;
  margin-right: 20px;
  cursor: pointer;
}
</style>