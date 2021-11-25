<template>
  <div class="main">
    <div class="login_main">
      <h3>登录</h3>
      <el-form ref="loginForm" class="login_form" :model="loginForm">
        <el-form-item>
          <el-input
            placeholder="用户名"
            size="medium"
            prefix-icon="el-icon-user"
            v-model="loginForm.username"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-input
            placeholder="密码"
            size="medium"
            type="password"
            prefix-icon="el-icon-lock"
            v-model="loginForm.password"
          ></el-input>
        </el-form-item>
      </el-form>
      <el-button
        type="primary"
        @click="handleLogin"
        :loading="loading"
        class="login_btn"
        >登录</el-button
      >
    </div>
  </div>
</template>
<script>
import { login } from '@/api/login.js'

export default {
  name: 'Login',
  data() {
    return {
      loginForm: {
        username: '',
        password: ''
      },
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      this.loading = true
      const { data } = await login({
        username: this.loginForm.username,
        password: this.loginForm.password
      })
      if (data.errCode !== 200) return this.$notify.error('登录失败！')
      window.localStorage.setItem('token', data.token)
      this.$router.push('/home')
      this.loading = false
    }
  }
}
</script>
<style lang="less" scoped>
.main {
  height: 100%;
  background: #304156;

  .login_main {
    width: 350px;
    height: 300px;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    line-height: 50px;
    text-align: center;
    padding: 20px;
    color: #fff;
    box-shadow: 0 0 0 200px rgba(255, 255, 255, 0.2) inset;
    border-radius: 15px;
    .login_btn {
      width: 100%;
    }
  }
}
</style>
