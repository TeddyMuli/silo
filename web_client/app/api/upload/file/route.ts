import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';
import { NextRequest, NextResponse } from 'next/server';
import path from 'path';
import mime from 'mime-types';
import { handleCreateFile } from '@/mutations';

export async function POST(request: NextRequest) {
  const s3 = new S3Client({
    endpoint: `https://${process.env.D_O_SPACES_URL}`,
    region: 'lon1',
    credentials: {
      accessKeyId: process.env.D_O_SPACES_ID!,
      secretAccessKey: process.env.D_O_SPACES_SECRET!,
    },
  });

  try {
    const formData = await request.formData();
    const file = formData.get('file') as File;
    const organization_id = formData.get('organization_id') as string;
    const folder_id = formData.get('folder_id') as string;

    if (!file) {
      return NextResponse.json({ error: 'No file uploaded' }, { status: 400 });
    }

    const fileContent = Buffer.from(await file.arrayBuffer());
    const fileKey = `uploads/${Date.now()}_${path.basename(file.name)}`;
    const contentType = mime.lookup(fileKey) || 'application/octet-stream';

    const command = new PutObjectCommand({
      Bucket: process.env.D_O_SPACES_URL!,
      Key: fileKey,
      Body: fileContent,
      ContentType: contentType,
      ACL: 'public-read',
    });

    const uploadResult = await s3.send(command);

    const fileToSend = {
      name: path.basename(file.name),
      file_path: `https://${process.env.D_O_SPACES_URL}/${process.env.D_O_SPACES_URL}/${fileKey}`,
      file_size: file.size,
      organization_id,
      folder_id,
    };

    const createFileResponse = await handleCreateFile(fileToSend);

    if (createFileResponse.status === 500) {
      return NextResponse.json({ error: 'Failed to save file metadata' }, { status: 500 });
    }
    
    return NextResponse.json({ message: 'File uploaded successfully', file: fileToSend.file_path }, { status: 200 });
  } catch (uploadError) {
    console.error('Upload error:', uploadError);
    return NextResponse.json({ error: 'File upload failed' }, { status: 500 });
  }
}
