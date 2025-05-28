"use client";
import { DownloadPopup } from "@/components/downloadPopup";
import { Navbar } from "@/components/navbar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import axios from "axios";
import { Clock, Scissors, Loader2 } from "lucide-react";
import React, { useState } from "react";
import { motion as m } from "motion/react";
import {
  validateInputs,
} from "@/lib/utils";

const BACKEND_ENDPOINT = process.env.NEXT_PUBLIC_BACKEND_URL;

const page = () => {
  const [url, setUrl] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [loading, setLoading] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [response, setResponse] = useState(null);

  const showTimestamps = url.length > 0;

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!validateInputs({ url, startTime, endTime })) {
      return;
    }
    setLoading(true);
    try {
      const response = await axios.get(
        `${BACKEND_ENDPOINT}/youtube/download?url=${url}&start=${startTime}&end=${endTime}`,
        { responseType: "blob" },
      );

      setResponse(response.data);
      setOpenDialog(true);
    } catch (error) {
      alert("Failed to download video. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Navbar />
      <DownloadPopup
        showDownloadModal={openDialog}
        setShowDownloadModal={setOpenDialog}
        response={response}
      />
      <div
        style={{
          backgroundImage: "url('/hero-bg.webp')",
          backgroundSize: "cover",
          backgroundRepeat: "no-repeat",
          backgroundPosition: "center",
        }}
        className="h-screen text-white w-full flex flex-col items-center justify-center"
      >
        <m.h1
          initial={{ opacity: 0, y: 20, filter: "blur(10px)" }}
          animate={{ opacity: 1, y: 0, filter: "blur(0px)" }}
          transition={{ duration: 0.5 }}
          className="md:text-4xl text-3xl mb-4 text-white text-center"
        >
          What do you want to clip?
        </m.h1>

        <div className="max-w-2xl mx-auto">
          <div className="p-6">
            <form onSubmit={handleSubmit} className="space-y-4">
              <m.div
                initial={{ opacity: 0, filter: "blur(10px)" }}
                animate={{ opacity: 1, filter: "blur(0px)" }}
                transition={{ duration: 0.8 }}
                className="space-y-2"
              >
                <Input
                  type="url"
                  placeholder="Paste YouTube URL here..."
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  className="h-12 text-base transition-all duration-200"
                  required
                  disabled={loading}
                />
              </m.div>

              <div
                className={`grid grid-cols-2 gap-4 transition-all duration-500 ease-out ${
                  showTimestamps
                    ? "opacity-100 max-h-20 translate-y-0"
                    : "opacity-0 max-h-0 -translate-y-4 overflow-hidden"
                }`}
              >
                <div className="space-y-1">
                  <label className="text-sm font-medium text-white flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    Start Time
                  </label>
                  <Input
                    type="text"
                    placeholder="0:00"
                    value={startTime}
                    onChange={(e) => setStartTime(e.target.value)}
                    className="h-10"
                    required={showTimestamps}
                    disabled={loading}
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-sm font-medium text-white flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    End Time
                  </label>
                  <Input
                    type="text"
                    placeholder="1:30"
                    value={endTime}
                    onChange={(e) => setEndTime(e.target.value)}
                    className="h-10"
                    required={showTimestamps}
                    disabled={loading}
                  />
                </div>
              </div>

              <m.div
                initial={{ opacity: 0, filter: "blur(10px)" }}
                animate={{ opacity: 1, filter: "blur(0px)" }}
                transition={{ duration: 1 }}
              >
                <Button
                  type="submit"
                  size="lg"
                  className="w-full h-12 text-base bg-gradient-to-r from-[#ccff00] to-[#b4ff00] hover:from-[#cfff2a] hover:to-[#ccff00] text-black transition-all duration-200 shadow-lg hover:shadow-xl"
                  disabled={!url || !startTime || !endTime || loading}
                >
                  {loading ? (
                    <span className="flex items-center justify-center">
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Clipping...
                    </span>
                  ) : (
                    <span className="flex items-center justify-center">
                      <Scissors className="w-4 h-4 mr-2" />
                      Create Clip
                    </span>
                  )}
                </Button>
              </m.div>
            </form>
          </div>
        </div>
      </div>
    </>
  );
};

export default page;
