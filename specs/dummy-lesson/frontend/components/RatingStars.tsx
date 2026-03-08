import React from 'react'

interface Props {
  value: number
}

export default function RatingStars({ value }: Props) {
  return <span>{'★'.repeat(value)}{'☆'.repeat(5 - value)}</span>
}
