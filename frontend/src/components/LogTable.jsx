// src/components/LogTable.js
import React from 'react';

function LogTable() {
  return (
    <div>
      <h2>日志展示区</h2>
      <table border="1">
        <thead>
          <tr><th>时间</th><th>IP</th><th>URI</th><th>动作</th></tr>
        </thead>
        <tbody>
          <tr><td>2025-06-20</td><td>192.168.1.1</td><td>/login</td><td>拦截</td></tr>
        </tbody>
      </table>
    </div>
  );
}

export default LogTable;

