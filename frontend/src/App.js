import React, { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [messages, setMessages] = useState([]);
  const [uploadedFileName, setUploadedFileName] = useState(null);

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    const MAX_FILE_SIZE = 10 * 1024 * 1024;

    if (selectedFile) {
      if (selectedFile.size > MAX_FILE_SIZE) {
        addMessage("File size exceeds the 10MB limit.", "bot");
        setFile(null);
        setUploadedFileName(null);
        return;
      }
      setFile(selectedFile);
      setUploadedFileName(selectedFile.name);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      addMessage("Please select a file before uploading.", "bot");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      addMessage("File uploaded successfully!", "bot");
      setFile(null);
      setUploadedFileName(null);
    } catch (error) {
      console.error("Error uploading file:", error);
      addMessage(error.response?.data || "Failed to upload the file.", "bot");
    }
  };

  const handleRemoveFile = () => {
    setFile(null);
    setUploadedFileName(null);
  };

  const handleChat = async () => {
    if (query.trim() === "") {
      addMessage("Please enter a query.", "bot");
      return;
    }

    addMessage(query, "user");

    try {
      const res = await axios.post("http://localhost:8080/chat", { query });

      if (res.data && res.data.answer) {
        addMessage(res.data.answer, "bot");
      } else if (res.data && res.data.error) {
        // Tangani jika server memberikan error tertentu
        addMessage(res.data.error, "bot");
      } else {
        addMessage("No valid response from the server.", "bot");
      }
    } catch (error) {
      console.error("Error querying chat:", error);
      const errorMessage =
        error.response?.data?.message ||
        error.message ||
        "Ada Sesuatu Yang Salah. Mohon Coba Lagi.";
      addMessage(errorMessage, "bot");
    }

    setQuery(""); // Reset input query setelah dikirim
  };

  const addMessage = (text, sender) => {
    setMessages((prevMessages) => [...prevMessages, { text, sender }]);
  };

  return (
    <div className="flex flex-col h-screen bg-gray-900">
      <div className="bg-gray-900 text-white text-center py-4">
        <h1 className="text-xl font-bold">Data Analysis Chatbot</h1>
      </div>

      <div className="flex-grow overflow-y-auto p-4 space-y-4">
        {messages.map((msg, index) => (
          <div
            key={index}
            className={`flex ${
              msg.sender === "user" ? "justify-end" : "justify-start"
            }`}
          >
            <div
              className={`max-w-[75%] p-3 rounded-lg text-sm ${
                msg.sender === "user"
                  ? "bg-blue-500 text-white"
                  : "bg-gray-200 text-gray-800"
              }`}
            >
              {msg.text}
            </div>
          </div>
        ))}
      </div>

      <div className="p-4 bg-grey-100 border-gray-100">
        {uploadedFileName && (
          <div className="mb-2 flex items-center justify-between text-gray-500">
            <p>
              File uploaded:{" "}
              <span className="font-semibold text-gray-800">
                {uploadedFileName}
              </span>
            </p>
            <button
              onClick={handleRemoveFile}
              className="text-red-500 hover:underline text-sm"
            >
              Remove
            </button>
          </div>
        )}

        <div className="flex items-center space-x-4">
          <input
            type="file"
            onChange={handleFileChange}
            className="hidden"
            id="file-input"
          />
          <label
            htmlFor="file-input"
            className="bg-gray-200 text-gray-600 py-2 px-4 rounded-lg cursor-pointer hover:bg-gray-300"
          >
            Upload
          </label>
          <button
            onClick={handleUpload}
            className="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600"
          >
            Analyze
          </button>

          <div className="flex items-center bg-gray-100 border border-gray-300 rounded-lg overflow-hidden w-full">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Type your question..."
              className="flex-grow bg-gray-100 py-2 px-4 text-gray-800 outline-none"
            />
            <button
              onClick={handleChat}
              className="bg-blue-500 text-white py-2 px-4 hover:bg-blue-600"
            >
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
