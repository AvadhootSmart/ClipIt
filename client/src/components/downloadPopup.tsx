import React, { useMemo, useEffect } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { Download } from "lucide-react";
import { Button } from "./ui/button";

interface Props {
  showDownloadModal: boolean;
  setShowDownloadModal: (value: boolean) => void;
  response: Blob | null;
}

export const DownloadPopup = ({
  showDownloadModal,
  setShowDownloadModal,
  response,
}: Props) => {
  const videoUrl = useMemo(() => {
    if (response) {
      return URL.createObjectURL(response);
    }
    return null;
  }, [response]);

  useEffect(() => {
    return () => {
      if (videoUrl) {
        URL.revokeObjectURL(videoUrl);
      }
    };
  }, [videoUrl]);

  const handleDownload = () => {
    if (!response) return;

    const link = document.createElement("a");
    link.href = videoUrl!;
    link.download = "ClipIt.mp4";
    link.click();
    setShowDownloadModal(false);
  };

  return (
    <Dialog open={showDownloadModal} onOpenChange={setShowDownloadModal}>
      <DialogContent className="min-w-[60vw] rounded-2xl bg-neutral-900 border-none text-white">
        <DialogHeader>
          <DialogTitle className="text-lg md:text-2xl text-left font-normal p-6">
            Your clip is <strong>ready!</strong>{" "}
          </DialogTitle>
        </DialogHeader>
        <div className="space-y-6 w-full">
          {videoUrl && (
            <div className="relative w-full">
              <video
                className="w-full aspect-video bg-muted object-cover focus:ring-0 focus:outline-none"
                poster="/placeholder.svg?height=200&width=350"
                controls
              >
                <source
                  src={videoUrl}
                  type="video/mp4"
                  className="w-full focus:ring-0"
                />
                Your browser does not support the video tag.
              </video>
            </div>
          )}

          <Button
            className="w-full h-12 text-base font-semibold bg-gradient-to-r font-normal text-black from-[#ccff00] to-[#aaff00] hover:from-[#bbff00] hover:to-[#99ee00] transition-all duration-200"
            onClick={handleDownload}
            disabled={!videoUrl}
          >
            <Download className="w-4 h-4 mr-2" />
            Download Clip
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
