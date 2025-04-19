import { useState, useRef, useEffect } from "react";
import { HiOutlineCommandLine } from "react-icons/hi2";
import styled from "styled-components";

import { useOpenai } from "../features/aicopilot/useOpenaiAsk";
import Button from "../ui/Button";
import Form from "../ui/Form.jsx";
import Row from "../ui/Row";
import Heading from "../ui/Heading";

const Textarea = styled.textarea`
  padding: 0.8rem 1.2rem;
  border: 1px solid var(--color-grey-300);
  border-radius: 5px;
  background-color: var(--color-grey-0);
  box-shadow: var(--shadow-sm);
  width: 100%;
  height: 8rem;
`;

const StyledButton = styled.div`
  &:has(button) {
    display: flex;
    justify-content: flex-end;
    align-items: flex-start;
    gap: 1.2rem;
  }
`;

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
          <StyledButton>
            <Heading as="h1">
              <HiOutlineCommandLine />{" "}
            </Heading>

            <Textarea
              type="text"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Type your message..."
            />
            <Button type="submit">Send</Button>
          </StyledButton>
        </Form>
      </Row>
    </>
  );
}

export default AICopilot;
