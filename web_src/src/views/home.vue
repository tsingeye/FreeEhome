<template>
  <el-container>
    <el-aside :width="isCollapse ? '64px' : '250px'">
      <el-menu
        :default-active="$route.path"
        unique-opened
        router
        background-color="#304156"
        :collapse="isCollapse"
        text-color="#ccc"
        :collapse-transition="false"
      >
        <!--el-menu-item index="/monitor">
          <i class="el-icon-menu"></i>
          <span slot="title">设备点播</span>
        </el-menu-item-->
        <el-menu-item index="/devices">
          <i class="el-icon-menu"></i>
          <span slot="title">设备管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="collapse" @click="isCollapse = !isCollapse">
          <i class="iconfont icon-hanbaocaidanzhedie"></i>
        </div>

        <div class="userInfo">
          <el-image
            class="avater"
            fit="cover"
            :src="userInfo.avater"
          ></el-image>
          <el-dropdown size="medium" @command="handleCommand">
            <span class="el-dropdown-link">
              {{ userInfo.name }}
              <i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="setting">设置</el-dropdown-item>
              <el-dropdown-item command="logout">注销</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <transition name="slide-fade" mode="out-in">
          <router-view />
        </transition>
      </el-main>
    </el-container>
  </el-container>
</template>
<script>
import { logout } from '@/api/login.js'
export default {
  name: 'Home',
  data() {
    return {
      userInfo: {
        name: 'superAdmin',
        avater: 'https://z3.ax1x.com/2021/11/24/oPDJk4.jpg'
      },
      isCollapse: false
    }
  },
  methods: {
    async handleCommand(command) {
      if (command === 'setting') {
        console.log('setting')
      } else if (command === 'logout') {
        const token = window.localStorage.getItem('token')
        const { data } = await logout(token)
        if (data.errCode !== 200) return this.$notify.error('请求失败！')
        this.$router.push('/')
      }
    }
  }
}
</script>
<style lang="less" scoped>
.el-container {
  height: 100%;
  .el-menu {
    height: 100%;
  }
}

.el-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  .collapse {
    letter-spacing: 0.2em;
    cursor: pointer;
    .icon-hanbaocaidanzhedie {
      font-size: 30px;
    }
  }
  .userInfo {
    display: flex;
    align-items: center;
    justify-content: center;
    .avater {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      margin-right: 5px;
    }
  }
}
</style>
