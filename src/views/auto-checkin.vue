<template>
  <el-container direction="vertical">
    <el-row justify="center" type="flex">
      <el-header>
        <h1 align="center" style="font-size: 24px">
          SCU健康每日报自动打卡工作流
        </h1>
        <div align="center">
          <img
            align="center"
            alt="runtime"
            hspace="5"
            src="https://img.shields.io/endpoint?url=https%3A%2F%2Fwww.scubot.com%2F.netlify%2Ffunctions%2Fbadge%3Ftype_%3D0"
          />
          <img
            align="center"
            alt="total"
            hspace="5"
            src="https://img.shields.io/endpoint?url=https%3A%2F%2Fwww.scubot.com%2F.netlify%2Ffunctions%2Fbadge%3Ftype_%3D1"
          />
          <img
            align="center"
            alt="valid workflow"
            hspace="5"
            src="https://img.shields.io/endpoint?url=https%3A%2F%2Fwww.scubot.com%2F.netlify%2Ffunctions%2Fbadge%3Ftype_%3D2"
          />
        </div>
      </el-header>
    </el-row>
    <br />
    <div
      align="center"
      style="
        display: flex;
        display: -webkit-flex;
        text-align: center;
        justify-content: center;
        align-items: center;
      "
    >
      <el-carousel :interval="4000" height="80px" style="width: 600px">
        <el-carousel-item v-for="item in display_list" :key="item">
          <h2 class="medium">{{ item }}</h2>
        </el-carousel-item>
      </el-carousel>
    </div>
    <el-row justify="center" type="flex">
      <el-col :span="17">
        <el-form ref="scu" :model="scu" label-width="115px">
          <el-form-item label="常用浏览器UA">
            <el-input v-model="scu.userAgent"></el-input>
          </el-form-item>

          <el-form-item label="微服务Cookies">
            <el-input
              v-model="scu.cookies"
              placeholder="
使用Chrome或Edge登录https://wfw.scu.edu.cn/ncov/wap/default/index（健康每日报网页版）
F12，切换至Network选项卡，刷新页面
将index请求的Cookies复制到此处
"
              :rows="5"
              type="textarea"
              @focus="showCookiesNotify"
              @blur="closeCookiesNotify"
            ></el-input>
            <el-link
              type="primary"
              href="https://wfw.scu.edu.cn/ncov/wap/default/index"
              target="_blank"
              >健康每日报网页链接</el-link
            >
          </el-form-item>
          <el-form-item label="每日打卡时间">
            <el-time-select
              @focus="showTimeHint"
              @blur="closeTimeHint"
              v-model="scu.triggerTime"
              :picker-options="{
                start: '00:10',
                step: '00:30',
                end: '13:30',
              }"
              placeholder="选择时间"
            >
            </el-time-select>
          </el-form-item>
          <el-form-item label="每日打卡地点">
            <el-input v-model="scu.area" @click.native="onMapSelect">
              <i slot="prefix" class="el-input__icon el-icon-map-location"></i>
            </el-input>
          </el-form-item>
          <el-form-item label="接受消息的QQ">
            <el-row>
              <el-col :span="1">
                <el-avatar :size="36" src="https://empty.com">
                  <img :src="avatar_url" />
                </el-avatar>
              </el-col>
              <el-col :span="23" style="padding-left: 10px">
                <el-input
                  v-model="scu.qqid"
                  placeholder="填写自己的QQ用于接受打卡回执，留空为不接受回执。"
                  @change="qqChanged"
                >
                </el-input>
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item v-show="false" label="Access Token">
            <el-input
              v-model="accessToken"
              show-password
              placeholder="使用 [POST] 方法请求 https://scubot.com/.netlify/functions/token 获取，表单数据置空即可"
            ></el-input>
          </el-form-item>
          <div class="scumap" align="center">
            <el-dialog
              title="选择定位"
              :visible.sync="dialogMapVisible"
              center
              top="15px"
            >
              <scumap ref="scumap" />
              <span slot="footer" class="dialog-footer">
                <el-button @click="dialogMapVisible = false">取 消</el-button>
                <el-button type="primary" @click="confirmLoc">确 定</el-button>
              </span>
            </el-dialog>
            <el-dialog
              title="如何获取Cookies"
              :visible.sync="dialogVideoVisible"
              center
              top="10px"
              width="80%"
            >
              <!-- <cookies-video /> -->
              <bili-video
                ihtml="https://player.bilibili.com/player.html?aid=633770199&bvid=BV1Rb4y1h7TQ&cid=432741794&page=1&danmaku=0&high_quality=1"
              />
              <!-- <bili-video ihtml="https://player.bilibili.com/player.html?aid=90992146&cid=155380603&page=1&danmaku=0&high_quality=1" /> -->
            </el-dialog>
            <el-dialog
              title="日志汇总"
              :visible.sync="dialogLogVisible"
              center
              top="10px"
              width="70%"
            >
              <Log :uid="this.preview.uid" :logs="this.logs" />
            </el-dialog>
          </div>
          <br />
          <el-form-item>
            <el-button
              type="warning"
              @click="onPreview"
              style="width: 100px"
              :loading="loading_check"
              icon="el-icon-circle-check"
              >验证</el-button
            >
            <el-button
              type="primary"
              @click="onSubmit"
              style="width: 100px; margin-left: 20px"
              :loading="loading_submit"
              icon="el-icon-upload2"
            >
              提交
            </el-button>
            <transition name="el-zoom-in-center">
              <el-button
                v-show="!newClient"
                type="success"
                @click="onLog"
                style="width: 100px; margin-left: 20px"
                icon="el-icon-document"
              >
                日志
              </el-button>
            </transition>
            <transition name="el-zoom-in-center">
              <a-popconfirm
                title="确定要停止自动打卡吗？"
                ok-text="确定"
                cancel-text="返回"
                @confirm="onDelete"
              >
                <el-button
                  v-show="!newClient"
                  type="danger"
                  style="width: 100px; margin-left: 20px"
                  icon="el-icon-delete"
                  :loading="loading_delete"
                >
                  删除
                </el-button>
              </a-popconfirm>
            </transition>
          </el-form-item>
          <el-form-item>
            <json-viewer
              :value="preview"
              :expand-depth="5"
              copyable
              boxed
              sort
              expanded
            ></json-viewer>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </el-container>
</template>


<script>
import axios from "axios";
import scumap from "../components/scumap.vue";
import BiliVideo from "../components/bili-video.vue";
import Log from "../components/log.vue";

function sleep(time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

export default {
  components: { scumap, BiliVideo, Log },
  data() {
    return {
      scu: {
        cookies: "",
        userAgent: "你的浏览器UA",
        triggerTime: "00:10",
        qqid: "",
        area: "四川省 成都市 武侯区",
        uid: "",
      },
      accessToken: "",
      preview: "请点击'验证'",
      dialogMapVisible: false,
      dialogVideoVisible: false,
      dialogLogVisible: false,
      cacheLocation: {
        lat: 30.630839301216,
        lng: 104.079966362848,
      },
      avatar_url: "https://q1.qlogo.cn/g?b=qq&nk=1&s=640",
      loading_check: false,
      loading_submit: false,
      loading_delete: false,
      display_list: [
        "成都服务器集群，不固定出口IP",
        "自定义打卡定位地点，支持范围选择",
        "保存原有个人信息模板，只针对定位信息做更改",
        "QQ机器人准确反馈打卡状态",
        "在线查询打卡日志(TODO)",
      ],
      newClient: true,
      logs: [],
    };
  },
  created() {
    this.scu.userAgent = navigator.userAgent;
    this.$notify.success({
      title: "提示",
      dangerouslyUseHTMLString: true,
      message:
        "如果你需要每日打卡消息回执，请主动添加QQ(机器人)：<strong><a href='https://wpa.qq.com/msgrd?v=3&uin=3583618673&site=scubot&from=scubot&menu=yes'>3583618673</a></strong>，验证消息填：<strong><span id='verifyCode' onclick='copy();'>7355608</span></strong>",
      duration: 0,
    });
    sleep(500).then(() => {
      this.$notify.warning({
        title: "使用须知",
        dangerouslyUseHTMLString: true,
        message:
          "在使用本项目提供的服务前，请确保知悉<strong><a target='_blank' href='/pdf/guideline.pdf'>项目导览</a></strong>中的内容",
        duration: 0,
      });
    });
  },
  methods: {
    onPreview() {
      this.loading_check = true;
      if (this.scu.cookies == null) {
        this.$message.error("<微服务Cookies>不能为空");
        this.loading_check = false;
        return;
      }
      const ePattern = /eai-sess=(.+?)(?:;|$| )/g;
      const uPattern = /UUkey=(.+?)(?:;|$| )/g;
      const eGot = ePattern.exec(this.scu.cookies);
      const uGot = uPattern.exec(this.scu.cookies);
      var cookies = {};
      if (eGot != null && eGot.length > 1 && uGot != null && uGot.length > 1) {
        cookies = {
          "eai-sess": eGot[1],
          UUkey: uGot[1],
        };
        axios
          .get("https://www.scubot.com/.netlify/functions/uuid", {
            params: {
              UUkey: cookies["UUkey"],
              "eai-sess": cookies["eai-sess"],
            },
          })
          .then((resp) => {
            this.scu.uid = resp.data["uid"];
            this.newClient = !resp.data["exist"];
            this.$message.success("cookies验证成功");

            this.preview = JSON.parse(JSON.stringify(this.scu));
            delete this.preview.area;
            this.preview.cookies = cookies;
            this.preview.location = this.cacheLocation;
            this.loading_check = false;
          })
          .catch((error) => {
            console.log(error);
            this.$message.error("出现了一些错误");
            this.loading_check = false;
          });
      } else {
        console.log(eGot, uGot);
        this.$message.error("<微服务Cookies>格式错误");
        this.loading_check = false;
        return;
      }
    },
    onSubmit() {
      if (this.preview == "请点击'验证'") {
        this.$message.warning("请先点击验证，再点击提交");
        return;
      }
      this.loading_submit = true;
      const postData = {
        content: JSON.parse(JSON.stringify(this.preview)),
        accessToken: this.accessToken,
        uid: this.preview.uid,
      };
      axios
        .post("https://www.scubot.com/.netlify/functions/post", postData)
        .then((res) => {
          if (res.status == 200) {
            this.$alert(res.data["message"], {
              confirmButtonText: '确定'
            });
          } else if (res.status == 403) {
            this.$message.error("错误：" + res.data["detail"]);
          } else {
            this.$message.error("出现了一些错误，请检查表单是否填写正确");
          }
          this.loading_submit = false;
        })
        .catch((res) => {
          this.$message.error("出现了一些错误，请检查表单是否填写正确");
          console.log(res);
          if (res.response != undefined && res.response.status == 403) {
            sleep(500).then(() => {
              this.$message.error("错误：" + res.response.data["detail"]);
            });
          }
          this.loading_submit = false;
        });
    },
    onLog() {
      axios
        .get("https://www.scubot.com/.netlify/functions/log", {
          params: {
            uid: this.preview.uid,
          },
        })
        .then((resp) => {
            this.logs = resp.data['message'];
            this.$message.success('日志获取成功');
            console.log(this.logs);
        })
        .catch((err) => {
            this.$message.error("日志获取失败");
            console.log(err);
        })
      this.dialogLogVisible = true;
    },
    onDelete() {
      this.loading_delete = true;
      const postData = {
        accessToken: this.accessToken,
        uid: this.preview.uid,
      };
      axios
        .post("https://www.scubot.com/.netlify/functions/delete", postData)
        .then((res) => {
          if (res.status == 200) {
            this.$message.success(res.data["message"]);
          } else if (res.status == 403) {
            this.$message.error("错误：" + res.data["detail"]);
          } else {
            this.$message.error("出现了一些错误，请检查表单是否填写正确");
          }
          this.loading_delete = false;
        })
        .catch((res) => {
          if (res.response != undefined && res.response.status == 403) {
            sleep(500).then(() => {
              this.$message.error("错误：" + res.response.data["detail"]);
            });
          }
          this.loading_delete = false;
          this.$message.error("出现了一些错误，请检查表单是否填写正确");
        });
    },
    onMapSelect() {
      this.dialogMapVisible = true;
    },
    confirmLoc() {
      this.dialogMapVisible = false;
      this.cacheLocation["lat"] = this.$refs.scumap.position[1];
      this.cacheLocation["lng"] = this.$refs.scumap.position[0];
      // get new area info
      this.$jsonp("https://api.map.baidu.com/reverse_geocoding/v3/", {
        output: "json",
        coordtype: "wgs84ll",
        ak: "0hYGiH3Ob5ZhV0eWzrGVXCD3bEdBCi6L",
        location: `${this.cacheLocation["lat"]},${this.cacheLocation["lng"]}`,
      })
        .then((response) => {
          if (response.status != 0) {
            console.log(response);
            this.$message.error("获取地理位置信息失败");
          } else {
            console.log(response.result);
            var addr = response.result.addressComponent;
            this.scu.area = `${addr.province} ${addr.city} ${addr.district}`;
            this.$message.success("地理位置获取：" + this.scu.area);
          }
        })
        .catch((error) => {
          console.log(error.response);
          this.$message.error("获取地理位置信息失败");
        });
    },
    showCookiesNotify() {
      this.cookies_not_inst = this.$notify({
        title: "如何获取cookies",
        dangerouslyUseHTMLString: true,
        duration: 0,
        onClick: this.cookiesNotifyClicked,
        showClose: false,
        message: "点击这里，可观看教程视频<i class='el-icon-video-play'/>",
      });
    },
    closeCookiesNotify() {
      this.cookies_not_inst.close();
    },
    cookiesNotifyClicked() {
      this.dialogVideoVisible = true;
    },
    qqChanged(value) {
      this.avatar_url = `https://q1.qlogo.cn/g?b=qq&nk=${value}&s=640`;
    },
    showTimeHint() {
      this.timeHintInst = this.$notify.info({
        title: "触发时间的随机性",
        showClose: false,
        dangerouslyUseHTMLString: true,
        message: "实际触发时间会包含<strong>正负1分钟内</strong>的随机偏移",
        duration: 0,
      });
    },
    closeTimeHint() {
      this.timeHintInst.close();
    },
  },
};
</script>
<style scoped>
/* @import url(https://unpkg.com/element-ui/lib/theme-chalk/index.css); */
.el-container {
  padding: 10px;
  font-family: Helvetica;
  margin-top: 10px;
}

.el-form {
  margin-top: 50px;
}

.scumap {
  width: 100%;
  height: 100%;
}
.el-carousel__item h2 {
  color: #2a3344;
  font-size: 20px;
  opacity: 0.75;
  line-height: 100px;
}

.el-carousel__item:nth-child(2n) {
  background-color: #fff;
}

.el-carousel__item:nth-child(2n + 1) {
  background-color: #fff;
}
</style>
