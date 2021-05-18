<template>
  <div class="login_out">
    <div class="login_container">
      <div class="login_title">FreeEhome</div>
      <el-form
        :label-position="'right'"
        label-width="80px"
        :model="ruleForm"
        ref="ruleForm"
        :rules="rules"
      >
        <el-form-item
          label="用户名"
          prop='username'
        >
          <el-input v-model="ruleForm.username"></el-input>
        </el-form-item>
        <el-form-item
          label="密码"
          prop='password'
        >
          <el-input v-model="ruleForm.password"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="submitForm('ruleForm')"
          >登录</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>
<script>
import axios from '@/utils/request'

export default {
  data () {
    return {
      labelPosition: 'right',
      ruleForm: {
        username: '',
        password: '',
      },
      rules: {
        username: { required: true, message: '请输入用户名', trigger: 'blur' },
        password: { required: true, message: '请输入密码', trigger: 'blur' },
      }
    };
  },
  methods: {
    submitForm (formName) {
      const that = this
      this.$refs[formName].validate((valid) => {
        if (valid) {
          axios.get('/api/v1/system/login', that.ruleForm)
            .then(res => {
              localStorage.setItem('authCode', res.authCode)
              that.$router.push('/')
            })
          //   alert('submit!');
        } else {
          console.log('error submit!!');
          return false;
        }
      });
    }
  }
}
</script>
<style lang="scss">
.login_out {
  height: 100vh;
  display: flex;
  align-items: center;
  background: linear-gradient(-135deg, #c850c0, #4158d0);
}
.login_container {
  width: 500px;
  margin: auto;
  padding: 20px;
  background: #fff;
  border-radius: 20px;
}
.login_title {
  margin-bottom: 50px;
  text-align: center;
}
</style>