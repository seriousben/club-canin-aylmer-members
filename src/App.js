// src/App.js
import React from "react";
import LazyLoad from "react-lazyload";
import { Avatar, Timeline, Layout, Tooltip } from "antd";
import groupBy from "lodash-es/groupBy";
import map from "lodash-es/map";
import "./App.css";

const { Header, Content } = Layout;

function formatDate(date) {
  var monthNames = [
    "Janvier/January",
    "Fevrier/February",
    "Mars/March",
    "Avril/April",
    "Mai/May",
    "Juin/June",
    "Juillet/July",
    "Aout/August",
    "Septembre/September",
    "Octobre/October",
    "Novembre/November",
    "Decembre/December"
  ];

  var day = date.getDate();
  var monthIndex = date.getMonth();
  var year = date.getFullYear();

  return day + " " + monthNames[monthIndex] + " " + year;
}

function App() {
  const [members, setMembers] = React.useState([]);

  React.useEffect(() => {}, []);

  React.useEffect(() => {
    async function fetchData() {
      const res = await fetch("/members.json", {});
      const mJson = await res.json();
      const ms = mJson.members.map(m => {
        m.imageUploadedAt = new Date(m.imageUploadedAt);
        return m;
      });
      ms.sort((a, b) => a.imageUploadedAt - b.imageUploadedAt);
      setMembers(ms);
    }
    fetchData();
  }, []);

  return (
    <Layout>
      <Header>
        <h1 style={{ color: "white" }}>Club Canin Aylmer</h1>
      </Header>
      <Content
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          marginTop: "25px"
        }}
      >
        <h2 style={{ marginBottom: "25px" }}>Chronologie / Timeline</h2>
        <Timeline mode="alternate" style={{ minWidth: "80%" }}>
          {map(
            groupBy(members, m => formatDate(m.imageUploadedAt)),
            (ms, i) => {
              return (
                <Timeline.Item key={i}>
                  <h3>{i}</h3>
                  <div>
                    {ms.map((m, j) => (
                      <LazyLoad height={64}>
                        <Tooltip key={j} title={m.name}>
                          <Avatar
                            onClick={() => window.open(m.imageUrl, "_blank")}
                            shape="square"
                            size={64}
                            src={m.imageUrl.replace("_orig", "")}
                          />
                        </Tooltip>
                      </LazyLoad>
                    ))}
                  </div>
                </Timeline.Item>
              );
            }
          )}
        </Timeline>
      </Content>
    </Layout>
  );
}

export default App;
