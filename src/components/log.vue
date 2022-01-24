<template>
  <div>
    <a-table :columns="this.columns" :data-source="this.clogs" rowKey="date">
      <span slot="date"><a-icon type="calendar-o" /> 日期 </span>
      <span slot="time"><a-icon type="clock-circle" /> 时间 </span>
      <span slot="area"><a-icon type="environment" /> 地点 </span>

      <span slot="ip" slot-scope="ip">
        {{ ip }}
      </span>
      <span slot="result" slot-scope="text, record">
        <a-tag
          :color="
            record.status > 0
              ? record.status === 1
                ? 'geekblue'
                : 'volcano'
              : '#87d068'
          "
        >
          {{ record.result }}
        </a-tag>
        <a-divider type="vertical" />
        <a-icon :type="record.notified ? 'check' : 'close'" />
      </span>
    </a-table>
  </div>
</template>
<script>
export default {
  props: ["uid", "logs"],
  data() {
    return {
      columns: [
        {
          dataIndex: "date",
          key: "date",
          slots: { title: "date" },
          scopedSlots: { customRender: "name" },
        },
        {
          dataIndex: "time",
          key: "time",
          slots: { title: "time" },
        },
        {
          dataIndex: "area",
          key: "area",
          slots: { title: "area" },
        },
        {
          title: "IP",
          key: "ip",
          dataIndex: "ip",
          scopedSlots: { customRender: "ip" },
        },
        {
          title: "执行结果 | 通知",
          key: "result",
          scopedSlots: { customRender: "result" },
        },
      ],
      clogs: this.logs,
      cuid: this.uid,
    };
  },
  watch: {
    uid(newVal) {
      this.cuid = newVal;
    },
    logs(newVal) {
      this.clogs = newVal;
    },
  },
  methods: {},
};
</script>
