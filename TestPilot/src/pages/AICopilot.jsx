import { useState, useRef, useEffect } from "react";
import { HiOutlineCommandLine } from "react-icons/hi2";

import { useOpenai } from "../features/aicopilot/useOpenaiAsk";
import Button from "../ui/Button";
import Form from "../ui/Form.jsx";
import Row from "../ui/Row";
import Heading from "../ui/Heading";
// import Textarea from "../ui/Textarea";
import Input from "../ui/Input";

function AICopilot() {
  const [inputText, setInputText] = useState("");
  const [messages, setMessages] = useState([]);
  const chatContainerRef = useRef(null);

  const { askAI } = useOpenai();

  useEffect(() => {
    // 滚动到聊天底部
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop =
        chatContainerRef.current.scrollHeight;
    }
  }, [messages]); // 当 messages 变化时滚动

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!inputText.trim()) return; // 阻止发送空消息

    // 添加用户消息到聊天列表
    setMessages((prevMessages) => [
      ...prevMessages,
      { text: inputText, isUser: true },
    ]);
    setInputText(""); // 清空输入框

    try {
      //使用createAsk发起 api 请求
      askAI(
        { userInput: inputText },
        {
          onSuccess: (reply) => {
            // 添加 ChatGPT 的回复到聊天列表
            const data = reply?.data?.data;
            console.log("reply 的数据 ", reply);
            setMessages((prevMessages) => [
              ...prevMessages,
              {
                text: data,
                isUser: false,
              },
            ]);
          },
          error: (error) => {
            // 添加错误消息到聊天列表
            setMessages((prevMessages) => [
              ...prevMessages,
              {
                text: `Error: ${error.message}`,
                isUser: false,
              },
            ]);
          },
        }
      );
    } catch (error) {
      console.error("Error calling OpenAI API:", error);
      // 添加错误消息到聊天列表
      setMessages((prevMessages) => [
        ...prevMessages,
        { text: `Error: ${error.message}`, isUser: false },
      ]);
    }
  };

  // const handleKeyDown = (event) => {
  //   if (event.key === "Enter" && !event.shiftKey) {
  //     event.preventDefault(); // 阻止默认的换行行为
  //     handleSubmit();
  //   }
  // };

  return (
    <>
      <Heading as="h1">Ask ChatGPT</Heading>
      <Row>
        <div className="chat-container" ref={chatContainerRef}>
          {messages.map((message, index) => (
            <div
              key={index}
              className={`message ${
                message.isUser ? "user-message" : "bot-message"
              }`}
            >
              {message.text}
            </div>
          ))}
        </div>
      </Row>

      <Row>
        <Form type="modal" onSubmit={handleSubmit} className="input-form">
          <Heading as="h2">
            <HiOutlineCommandLine />{" "}
          </Heading>
          <Input
            type="text"
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            // onKeyDown={handleKeyDown}
            placeholder="Type your message..."
          />
          <Button type="submit">Send</Button>
        </Form>
      </Row>
    </>
  );
}

export default AICopilot;
