import React from 'react'

interface Props {
  url: string
}

export default function VideoPlayer({ url }: Props) {
  return <video src={url} controls className="w-full" />
}
